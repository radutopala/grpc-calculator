package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"

	calculator "github.com/radutopala/grpc-calculator/api"
	"github.com/radutopala/grpc-calculator/service"
	jaeger_metrics "github.com/uber/jaeger-lib/metrics"

	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	prometheus_metrics "github.com/uber/jaeger-lib/metrics/prometheus"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Panic handler prints the stack trace when recovering from a panic.
var panicHandler = grpc_recovery.RecoveryHandlerFunc(func(p interface{}) error {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	log.Errorf("panic recovered: %+v", string(buf))
	return status.Errorf(codes.Internal, "%s", p)
})

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Calculator"
	app.Version = "0.0.1"
	app.Flags = flags
	app.Action = start

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) {
	lis, err := net.Listen("tcp", c.String("bind-grpc"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", c.String("bind-grpc"))
	}

	// Logrus
	logger := log.NewEntry(log.New())
	grpc_logrus.ReplaceGrpcLogger(logger)
	log.SetLevel(log.InfoLevel)

	// Prometheus monitoring
	metrics := prometheus_metrics.New()

	// Jaeger tracing
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: c.Float64("jaeger-sampler"),
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: c.String("jaeger-host") + ":" + c.String("jaeger-port"),
		},
	}
	tracer, closer, err := cfg.New(
		"calculator",
		config.Logger(jaegerLoggerAdapter{logger}),
		config.Observer(rpcmetrics.NewObserver(metrics.Namespace(jaeger_metrics.NSOptions{Name: "calculator"}), rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Fatalf("Cannot initialize Jaeger Tracer %s", err)
	}
	defer closer.Close()

	// Set GRPC Interceptors
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_logrus.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
	)

	// Register Calculator service, prometheus and HTTP service handler
	calculator.RegisterServiceServer(server, &service.Service{})
	reflection.Register(server)
	grpc_prometheus.Register(server)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(c.String("bind-prometheus-http"), mux)
	}()

	log.Println("Starting Calculator service..")
	go server.Serve(lis)

	conn, err := grpc.Dial(c.String("bind-grpc"), grpc.WithInsecure())
	if err != nil {
		panic("Couldn't contact grpc server")
	}

	mux := grpc_runtime.NewServeMux()
	err = calculator.RegisterServiceHandler(context.Background(), mux, conn)
	if err != nil {
		panic("Cannot serve http api")
	}
	http.ListenAndServe(c.String("bind-http"), mux)
}

type jaegerLoggerAdapter struct {
	logger *log.Entry
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "bind-http",
		Usage:  "bind address for HTTP",
		EnvVar: "BIND_HTTP",
		Value:  ":8080",
	},
	cli.StringFlag{
		Name:   "bind-grpc",
		Usage:  "bind address for gRPC",
		EnvVar: "BIND_GRPC",
		Value:  ":2338",
	},
	cli.StringFlag{
		Name:   "bind-prometheus-http",
		Usage:  "bind prometheus address for HTTP",
		EnvVar: "BIND_PROMETHEUS_HTTP",
		Value:  ":8081",
	},
	cli.StringFlag{
		Name:   "jaeger-host",
		Usage:  "Jaeger hostname",
		EnvVar: "JAEGER_HOST",
		Value:  "127.0.0.1",
	},
	cli.IntFlag{
		Name:   "jaeger-port",
		Usage:  "Jaeger port",
		EnvVar: "JAEGER_PORT",
		Value:  5775,
	},
	cli.Float64Flag{
		Name:   "jaeger-sampler",
		Usage:  "Jaeger sampler",
		EnvVar: "JAEGER_SAMPLER",
		Value:  0.05,
	},
	cli.StringFlag{
		Name:   "jaeger-tags",
		Usage:  "Jaeger tags",
		EnvVar: "JAEGER_TAGS",
		Value:  "calculator",
	},
}
