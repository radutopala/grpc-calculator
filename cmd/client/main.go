package main

import (
	"context"
	calculator "github.com/radutopala/grpc-calculator/api"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"log"
	"os"
	"path"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Calculator gRPC Client"
	app.Version = "0.0.1"
	app.Flags = flags
	app.Action = start

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) {
	conn, err := grpc.Dial(c.String("grpc-address"), grpc.WithInsecure())
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		panic("Couldn't contact grpc server")
	}

	client := calculator.NewServiceClient(conn)
	response, _ := client.Compute(context.Background(), &calculator.Request{Expression: c.Args()[0]})

	log.Printf("Result is: %v", response.Result)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "grpc-address",
		Usage:  "gRPC for address",
		EnvVar: "GRPC_ADDRESS",
		Value:  ":2338",
	},
}
