package main

import (
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"net/url"
	"os"
	"strings"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	const nameEnvKey = "BEN_NAME"
	const subsystemEnvKey = "BEN_SUBSYSTEM"
	const serviceTypeEnvKey = "BEN_SERVICE_TYPE"
	const listenAddressEnvKey = "BEN_LISTEN_ADDRESS"
	const calleeEnvKey = "BEN_CALLEES"
	const calleeSeparator = " "

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	serviceTypeStr, exists := os.LookupEnv(serviceTypeEnvKey)
	serviceType := ServiceType(serviceTypeStr)
	if !exists {
		serviceType = vanilla
	}
	logger.Log("service_type", serviceType)

	listenAddress, exists := os.LookupEnv(listenAddressEnvKey)
	if !exists {
		listenAddress = ":8080"
	}
	logger.Log("listen_address", listenAddress)

	subsystem := os.Getenv(subsystemEnvKey)
	if subsystem == "" {
		subsystem = "subsystem"
	}
	name := os.Getenv(nameEnvKey)
	if name == "" {
		name = "name"
	}
	logger.Log("name", name, "subsystem", subsystem)

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "foo_metrics",
		Subsystem: "get_info",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, nil)
	requestLatency := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace:   "ben_base",
		Subsystem:   subsystem,
		Name:        name,
		Help:        "Total duration of requests in microseconds.",
		Buckets:     []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, nil)

	// Create callee endpoint list from the environment variable
	var callees []endpoint.Endpoint
	calleeList := os.Getenv(calleeEnvKey)
	logger.Log("callees", calleeList)
	for _, calleeStr := range strings.Split(calleeList, calleeSeparator) {
		callURL, err := url.Parse(calleeStr)
		if err != nil {
			panic(err)
		}
		calleeEndpoint := httptransport.NewClient(
			"GET",
			callURL,
			encodeRequest,
			decodeBaseResponse,
		).Endpoint()
		callees = append(callees, calleeEndpoint)
	}

	var svc BaseService
	svc = baseService{callees: callees, serviceType: serviceType}
	svc = loggingMiddleware(logger)(svc)
	svc = instrumentingMiddleware(requestCount, requestLatency)(svc)

	baseHandler := httptransport.NewServer(
		makeBaseEndPoint(svc),
		decodeBaseRequest,
		encodeResponse,
	)

	http.Handle("/", baseHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", listenAddress)
	logger.Log("err", http.ListenAndServe(listenAddress, nil))
}
