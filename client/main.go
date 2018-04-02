package main

import (
	"flag"
	"fmt"

	pb "github.com/igor-karpukhin/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
	"context"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:9090", "Server addr")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := pb.NewTestDataProviderClient(conn)

	resp, err := client.GetTestData(context.Background(), &pb.TestRequest{ID:0})
	if err != nil {
		panic(err)
	}

	fmt.Println("RESP ID:", resp.GetID())

	fmt.Println(*addr)
}
