package main

import (
	"github.com/go-kit/kit/endpoint"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	const nameEnvKey = "BEN_NAME"
	const subsystemEnvKey = "BEN_SUBSYSTEM"
	const serviceTypeEnvKey = "BEN_SERVICE_TYPE"
	const listenAddressEnvKey = "BEN_LISTEN_ADDRESS"
	const callEnvKey = "BEN_CALLS"
	const callSeparator = " "

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	serviceTypeStr, _ := getEnvString(serviceTypeEnvKey, string(vanilla))
	serviceType := ServiceType(serviceTypeStr)

	logger.Log("service_type", serviceType)

	listenAddress, _ := getEnvString(listenAddressEnvKey, ":8080")
	logger.Log("listen_address", listenAddress)

	subsystem, _ := getEnvString(subsystemEnvKey, "subsystem")
	name, _ := getEnvString(nameEnvKey, "name")
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

	// Create call endpoint list from the environment variable
	var calls []endpoint.Endpoint
	callList, exists := getEnvString(callEnvKey, "")
	if exists {
		logger.Log("calls", callList)

		for _, callStr := range strings.Split(callList, callSeparator) {
			callURL, err := url.Parse(callStr)
			if err != nil {
				panic(err)
			}
			callEndpoint := httptransport.NewClient(
				"GET",
				callURL,
				encodeRequest,
				decodeBaseResponse,
			).Endpoint()
			calls = append(calls, callEndpoint)
		}
	} else {
		logger.Log("calls", "[empty call list]")
	}

	var svc BaseService
	svc = baseService{calls: calls, serviceType: serviceType}
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

// Get an environment variable by name
// which is replaced by default value if it's empty
func getEnvString(key string, value string) (string, bool) {
	v := os.Getenv(key)
	if v == "" {
		return value, false
	}

	return v, true
}