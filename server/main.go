package main

import (
	"net"

	"context"
	"fmt"
	pb "github.com/albertocsm/grpc-client-balancing-test/grpc"
	"google.golang.org/grpc"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-metrics"
	"net/http"
	"time"
	"flag"
)

type TestServer struct {
	Data        *pb.TestResponse
	Ops         *prometheus.CounterVec
	OpsDuration *prometheus.HistogramVec
}

func (ts *TestServer) GetTestData(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {

	t := time.Now()
	defer func() {
		dt := time.Since(t)
		ts.countOperationTime("GetTestData", "GetTestData", dt)
	}()

	//fmt.Println("Request ID:", req.GetID())
	ts.Data.ID = req.GetID()
	return ts.Data, nil
}

func (ts *TestServer) countOperationTime(funcName string, operation string, dt time.Duration) {
	ts.Ops.WithLabelValues(funcName, operation).Inc()
	ts.OpsDuration.WithLabelValues(funcName, operation).Observe(dt.Seconds())
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
	metricsSystem := metricsCollector.NewSubsystem("gserver")
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// want to disregard loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {

	addr := flag.String("addr", "0.0.0.0:9090", "Bind addr")
	flag.Parse()

	opsCounter, opsHist := initMetrics()
	testServer := &TestServer{
		&pb.TestResponse{
			ID:      1,
			IntData: 100,
			StrData: "SomeData",
			IPAddr:  getLocalIP(),
		},
		opsCounter,
		opsHist,
	}

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("started to listen on", *addr)
	server := grpc.NewServer()
	pb.RegisterTestDataProviderServer(server, testServer)
	server.Serve(listener)
}
