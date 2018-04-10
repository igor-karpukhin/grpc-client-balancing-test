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
	t.Data.ID = req.GetID()
	return t.Data, nil
}

func main() {

	testServer := &TestServer{&pb.TestResponse{
		ID:      1,
		IntData: 100,
		StrData: "SomeData",
		IPAddr:  GetLocalIP(),
	}}

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

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
