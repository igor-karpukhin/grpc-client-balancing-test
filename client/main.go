package main

import (
	"flag"

	"github.com/go-kit/kit/transport/grpc"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9090", "Server addr")

	client := grpc.NewClient()
}
