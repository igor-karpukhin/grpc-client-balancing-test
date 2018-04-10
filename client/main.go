package main

import (
	"flag"
	"fmt"

	"context"
	pb "github.com/albertocsm/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"google.golang.org/grpc/naming"
)

func main() {

	addr := flag.String("addr", "0.0.0.0:9090", "Server addr")
	isDns := flag.Bool("dns", false, "Server addr")

	flag.Parse()

	var conn *grpc.ClientConn = nil

	if *isDns{

		fmt.Println("using DNS...")
		resolver, e := naming.NewDNSResolver()
		if e != nil {
			panic(e)
		}

		balancer := grpc.WithBalancer(grpc.RoundRobin(resolver))
		clientConn, err := grpc.Dial(*addr, grpc.WithInsecure(), balancer)
		if err != nil {
			panic(err)
		}

		conn = clientConn
	} else {

		fmt.Println("using target...")
		clientConn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		conn = clientConn
	}

	//rand.Seed(int64(time.Now().Second()))
	client := pb.NewTestDataProviderClient(conn)

	fmt.Println(fmt.Sprintf("connected to [%s]... ", *addr))

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case s := <-sig:

			fmt.Println("signal received: ", s)
		case _ = <-time.After(1 * time.Second):

			fmt.Println("request sent...")
			resp, err := client.GetTestData(context.TODO(), &pb.TestRequest{ID: int32(rand.Int())})
			if err != nil {
				panic(err)
			}
			fmt.Println("Response received", resp.ID, resp.IPAddr)
		}
	}
	fmt.Println(*addr)
}
