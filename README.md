# Calculator

[![Github Actions](https://github.com/radutopala/grpc-calculator/workflows/tests/badge.svg)](https://github.com/radutopala/grpc-calculator/actions)

A GRPC protobuf k8s calculator made in Go.

## Usage

```
go run cmd/main.go
```
```
NAME:
   main - Calculator

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --bind-http value             bind address for HTTP (default: ":8080") [$BIND_HTTP]
   --bind-grpc value             bind address for gRPC (default: ":2338") [$BIND_GRPC]
   --bind-prometheus-http value  bind prometheus address for HTTP (default: ":8081") [$BIND_PROMETHEUS_HTTP]
   --jaeger-host value           Jaeger hostname (default: "127.0.0.1") [$JAEGER_HOST]
   --jaeger-port value           Jaeger port (default: 5775) [$JAEGER_PORT]
   --jaeger-sampler value        Jaeger sampler (default: 0.05) [$JAEGER_SAMPLER]
   --jaeger-tags value           Jaeger tags (default: "calculator") [$JAEGER_TAGS]
   --help, -h                    show help
   --version, -v                 print the version
```

This will provide 3 endpoints:
 * GRPC at `localhost:2338`
 * HTTP at `http://localhost:8080`
 * Prometheus at `http://localhost:8081`

## Calls

### HTTP

```
curl -X POST -H "Content-Type: application/json" -d '{"expression":"3+5+(10*2)"}' "http://localhost:8080/compute"
{"result": "28"}
```

### GRPC

To install `grpc_cli` on a Mac, run `brew install grpc`.

```
grpc_cli call localhost:2338 Compute 'expression: "3+5+(10*2)"'
connecting to localhost:2338
result: "28"

Rpc succeeded with OK status
```
