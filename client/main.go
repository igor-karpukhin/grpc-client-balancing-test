package main

import (
	"flag"
	"fmt"

	"context"
	pb "github.com/albertocsm/grpc-client-balancing-test/grpc"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-metrics"
	"net/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

func initConnection(isDns *bool, addr *string) *grpc.ClientConn {
	if *isDns {

		fmt.Println("using DNS...")
		resolver, e := naming.NewDNSResolverWithFreq(time.Second * time.Duration(30))
		//resolver, e := naming.NewDNSResolver()
		if e != nil {
			panic(e)
		}

		balancer := grpc.WithBalancer(grpc.RoundRobin(resolver))
		clientConn, err := grpc.Dial(*addr, grpc.WithInsecure(), balancer)
		if err != nil {
			panic(err)
		}

		return clientConn
	} else {

		fmt.Println("using target...")
		clientConn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		return clientConn
	}
}

func initMetrics() (*prometheus.CounterVec, *prometheus.HistogramVec) {

	machineName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	metricsCollector := metrics.NewMetricCollector(
		prometheus.NewRegistry(),
		"travelaudience",
		prometheus.Labels{"server": machineName})
	http.Handle("/metrics", metricsCollector.Handler())
	metricsSystem := metricsCollector.NewSubsystem("gclient")
	go http.ListenAndServe("0.0.0.0:8081", nil)
	opsCounter, opsHist := newDurationMetric(metricsSystem, "operations", "function", "op")
	return opsCounter, opsHist
}

func newDurationMetric(
	metricSystem *metrics.Subsystem,
	metricName string,
	fields ...string) (*prometheus.CounterVec, *prometheus.HistogramVec) {

	metricOps := metricSystem.Vec(fmt.Sprintf("%s_ops", metricName), fields)
	metricOpsDuration := metricSystem.VecHistogram(fmt.Sprintf("%s_ops_duration", metricName), fields,
		[]float64{0.00001, 0.00005, 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.075, 0.1})

	return metricOps, metricOpsDuration
}

func countOperationTime(opsCounter *prometheus.CounterVec, opsHist *prometheus.HistogramVec, funcName string, operation string, dt time.Duration) {
	opsCounter.WithLabelValues(funcName, operation).Inc()
	opsHist.WithLabelValues(funcName, operation).Observe(dt.Seconds())
}

func makeRequest(client pb.TestDataProviderClient, opsCounter *prometheus.CounterVec, opsHist *prometheus.HistogramVec) (*pb.TestResponse, error) {

	t := time.Now()
	defer func() {
		dt := time.Since(t)
		countOperationTime(opsCounter, opsHist, "GetTestData", "GetTestData", dt)
	}()

	//fmt.Println("Sending request...")
	resp, err := client.GetTestData(context.TODO(), &pb.TestRequest{ID: int32(rand.Int())})
	//fmt.Println("Response received", resp.ID, resp.IPAddr)

	return resp, err
}

func main() {

	addr := flag.String("addr", "0.0.0.0:9090", "Server addr")
	isDns := flag.Bool("dns", false, "Should use DNS resolver")
	frequency := flag.Int("frequency", 0, "Request frequency")
	flag.Parse()

	// need telemetry subsys
	opsCounter, opsHist := initMetrics()

	// need a connection and a client
	conn := initConnection(isDns, addr)
	client := pb.NewTestDataProviderClient(conn)
	fmt.Println(fmt.Sprintf("connected to [%s]... ", *addr))

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGHUP)

	// all done and ready to start pounding the server
	for {
		select {
		case s := <-sig:
			fmt.Println("signal received: ", s)
		case _ = <-time.After(time.Duration(*frequency) * time.Millisecond):
			_, err := makeRequest(client, opsCounter, opsHist)
			if err != nil {
				fmt.Println(err)
				// want to back-off for a moment
				time.Sleep(1 * time.Second)
			}
		}
	}
	fmt.Println(*addr)
}
