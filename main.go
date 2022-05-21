package main

import (
	"context"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	// -------------------
	// Declare constants
	// -------------------
	const (
		nameEnvKey          = "BEN_NAME"
		subsystemEnvKey     = "BEN_SUBSYSTEM"
		listenAddressEnvKey = "BEN_LISTEN_ADDRESS"
		dbReadOpsEnvKey     = "BEN_DB_READ_OPS"
		dbWriteOpsEnvKey    = "BEN_DB_WRITE_OPS"
		dbUserEnvKey        = "BEN_DB_USER"
		dbPasswordEnvKey    = "BEN_DB_PASSWORD"
		dbCollectionEnvKey  = "BEN_DB_COLLECTION"
	)

	const databaseName = "data"
	const collectionName = "items"

	// -------------------
	// Init logging
	// -------------------
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// -------------------
	// Parse Environment variables
	// -------------------
	var (
		dbReadOps    int
		dbWriteOps   int
		dbUser       string
		dbPassword   string
		dbCollection string
	)
	dbReadOps, _ = getEnvInt(dbReadOpsEnvKey, 1)
	dbWriteOps, _ = getEnvInt(dbWriteOpsEnvKey, 1)
	dbUser, _ = getEnvString(dbUserEnvKey, "")
	dbPassword, _ = getEnvString(dbPasswordEnvKey, "")
	dbCollection, _ = getEnvString(dbCollectionEnvKey, "")

	logger.Log("db read ops", dbReadOps)
	logger.Log("db write ops", dbWriteOps)
	logger.Log("db user", dbUser)
	logger.Log("db password", dbPassword)
	logger.Log("db collection", dbCollection)

	listenAddress, _ := getEnvString(listenAddressEnvKey, ":8080")
	logger.Log("listen_address", listenAddress)

	subsystem, _ := getEnvString(subsystemEnvKey, "subsystem")
	name, _ := getEnvString(nameEnvKey, "name")
	logger.Log("name", name, "subsystem", subsystem)

	// -------------------
	// Init Prometheus counter & histogram
	// -------------------
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ben_base",
		Subsystem: subsystem,
		Name:      "request_count",
		Help:      "Number of requests received.",
		ConstLabels: map[string]string{
			"bensim_service_name": name,
		},
	}, nil)
	latencyCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ben_base",
		Subsystem: subsystem,
		Name:      "latency_counter",
		Help:      "Processing time taken of requests in seconds, as counter.",
		ConstLabels: map[string]string{
			"bensim_service_name": name,
		},
	}, nil)
	latencyHistogram := kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "ben_base",
		Subsystem: subsystem,
		Name:      "latency_histogram",
		Help:      "Processing time taken of requests in seconds, as histogram.",
		// TODO: determine appropriate buckets
		Buckets: []float64{.0002, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 25},
		ConstLabels: map[string]string{
			"bensim_service_name": name,
		},
	}, nil)
	throughput := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "ben_base",
		Subsystem: subsystem,
		Name:      "throughput",
		Help:      "Size of data transmitted in bytes.",
		ConstLabels: map[string]string{
			"bensim_service_name": name,
		},
	}, nil)

	// -------------------
	// Init database connection
	// -------------------

	// Connect to database and locate the collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	credential := options.Credential{
		Username: "root",
		Password: "password",
	}
	clientOpts := options.Client().
		ApplyURI("mongodb://localhost").
		SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	collection := client.Database(databaseName).Collection(collectionName)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// -------------------
	// Create & run service
	// -------------------
	var svc BaseService
	svc = baseService{
		dbCollection: collection,
		dbReadOps:    dbReadOps,
		dbWriteOps:   dbWriteOps,
	}
	svc = loggingMiddleware(logger)(svc)
	svc = instrumentingMiddleware(requestCount, latencyCounter, latencyHistogram, logger)(svc)

	baseHandler := httptransport.NewServer(
		makeBaseEndPoint(svc),
		decodeBaseRequest,
		encodeResponse,
		// Request throughput instrumentation
		httptransport.ServerFinalizer(func(ctx context.Context, code int, r *http.Request){
			responseSize := ctx.Value(httptransport.ContextKeyResponseSize).(int64)
			logger.Log("reponse_size", responseSize)
			throughput.Add(float64(responseSize))
		}),
	)

	http.Handle("/", baseHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", listenAddress)
	logger.Log("err", http.ListenAndServe(listenAddress, nil))
}