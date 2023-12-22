package main

import (
	"context"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/ayadi-mohamed/videos-microservice/jaeger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	httproutermiddleware "github.com/slok/go-http-metrics/middleware/httprouter"
)

const (
	metricsAddr = ":8000"
)
const version = "1.0.0"
var environment = os.Getenv("ENVIRONMENT")
var redisHost = os.Getenv("REDIS_HOST")
var redisPort = os.Getenv("REDIS_PORT")
var password = os.Getenv("PASSWORD")
var flaky = os.Getenv("FLAKY")

var ctx = context.Background()
var rdb redis.UniversalClient
var Logger, _ = zap.NewProduction()
var Sugar = Logger.Sugar()
var traceProvider = jaeger.NewJaegerTracerProvider()

func main() {
	r := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{redisHost + ":" + redisPort},
		DB:       0,
		Password: password,
	})
	rdb = r

	// Create our middleware.
	promMiddleware := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	router := httprouter.New()

	router.GET("/", httproutermiddleware.Handler("/", HandleHealthz, promMiddleware))
	router.GET("/:id", httproutermiddleware.Handler("/:id", HandleGetVideoById, promMiddleware))

	Sugar.Infow("Running...")
	// Serve our metrics.
	go func() {
		Sugar.Infof("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	log.Fatal(http.ListenAndServe(":10010", router))

}

func video(writer http.ResponseWriter, request *http.Request, p httprouter.Params, ctx context.Context, requestID uuid.UUID, ip string) (response string) {
	tr := traceProvider.Tracer("videos-ms-main-component")
	ctx, span := tr.Start(ctx, "Fetch Videos From DB")

	defer span.End()
	span.SetAttributes(attribute.Key("Function").String("GET VIDEO FROM DB"))
	span.SetAttributes(attribute.Key("UUID").String(requestID.String()))
	span.SetAttributes(attribute.Key("Client IP").String(ip))

	id := p.ByName("id")
	Sugar.Infof("Getting video with the ID: %v\n", id)

	videoData, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		Sugar.Infof("Item with the ID: %v is not found\n", id)
		return "{}"
	} else if err != nil {
		Sugar.Errorw("Error while fetching video from DB", err)
		span.RecordError(err)
		panic(err)
	} else {
		Sugar.Infof("Successfully fetched the video with ID %v from redis\n", id)
		return videoData
	}
}
