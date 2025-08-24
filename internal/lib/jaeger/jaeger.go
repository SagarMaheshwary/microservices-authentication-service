package jaeger

import (
	"context"
	"log"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Init(ctx context.Context, url string) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(url),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(constant.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown
}
