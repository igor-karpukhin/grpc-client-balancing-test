package main

import (
	"flag"
	"net"

	"context"
	"fmt"
	pb "github.com/albertocsm/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
)

type TestServer struct {
	Data *pb.TestResponse
}

func (t *TestServer) GetTestData(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	fmt.Println("REQ ID:", req.GetID())
	return t.Data, nil
}

func main() {

	testServer := &TestServer{
		Data: &pb.TestResponse{
			ID:      1,
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
	fmt.Println("started to listen on", *addr)
	server := grpc.NewServer()
	pb.RegisterTestDataProviderServer(server, testServer)
	server.Serve(listener)
}
