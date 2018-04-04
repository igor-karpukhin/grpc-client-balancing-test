package main

import (
	"flag"
	"fmt"

	pb "github.com/igor-karpukhin/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math/rand"
)

func main() {
	addr := flag.String("addr", "test-gserver:9090", "Server addr")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	rand.Seed(int64(time.Now().Second()))
	client := pb.NewTestDataProviderClient(conn)

	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case s := <-sig:
			fmt.Println("signal received: ", s)
		case <-time.After(1 * time.Second):
			fmt.Println("request sent...")
			resp, err := client.GetTestData(context.Background(), &pb.TestRequest{ID: int32(rand.Int())})
			if err != nil {
				panic(err)
			}
			fmt.Println("Response received", resp.ID)
		}
	}
	fmt.Println(*addr)
}
