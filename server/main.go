package main

import (
	"flag"
	"net"

	pb "github.com/igor-karpukhin/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
	"context"
)

type TestServer struct {
	Data *pb.TestResponse
}

func (t *TestServer) GetTestData(context.Context, *pb.TestRequest) (*pb.TestResponse, error) {
	return t.Data, nil
}

func main() {

	testServer := TestServer{
		Data: &pb.TestResponse{
			ID: 1,
			IntData: 100,
			StrData: "SomeData",
		},
	}

	addr := flag.String("addr", "0.0.0.0:9090", "Bind addr")
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterTestDataProviderServer(server, testServer)
	server.Serve(listener)
}
