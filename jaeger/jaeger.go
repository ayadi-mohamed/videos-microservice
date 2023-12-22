package jaeger

import (
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var TraceProvider *trace.TracerProvider
var jaegerEndpoint = os.Getenv("JAEGER_ENDPOINT")

func NewJaegerTracerProvider() *trace.TracerProvider {
	if TraceProvider == nil {
		initializeProvider()
	}
	return TraceProvider
}

func initializeProvider() {
	TraceProvider, _ = JeagerProvider(jaegerEndpoint)
}

func JeagerProvider(url string) (*tracesdk.TracerProvider, error) {
	service := "videos-ms"
	id := 1
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(service),
				attribute.Int("ID", id),
			),
		),
	)
	return tp, nil
}
