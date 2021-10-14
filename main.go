package main

import (
	"flag"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	const calleeEnvKey = "BEN_CALLEES"
	const calleeSeparator = " "

	serviceTypeStr, exists := os.LookupEnv(serviceTypeEnvKey)
	var serviceType ServiceType
	if !exists {
		serviceType = vanilla
	} else {
		serviceTypeInt, err := strconv.Atoi(serviceTypeStr)
		if err != nil {
			panic(err)
		}
		serviceType = ServiceType(serviceTypeInt)
	}

	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

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
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
