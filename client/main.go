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
)

func main() {
	addr := flag.String("addr", "0.0.0.0:9090", "Server addr")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	rand.Seed(int64(time.Now().Second()))
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
			fmt.Println("Response received", resp.ID)
		}
	}
	fmt.Println(*addr)
}
