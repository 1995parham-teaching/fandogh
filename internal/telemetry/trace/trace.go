package trace

import (
	"context"
	"log"

	telemetryConfig "github.com/1995parham-teaching/fandogh/internal/telemetry/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func New(cfg telemetryConfig.Trace) trace.Tracer {
	var exporter sdktrace.SpanExporter

	var err error
	if !cfg.Enabled {
		exporter, err = stdout.New(
			stdout.WithPrettyPrint(),
		)
	} else {
		exporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithEndpoint(cfg.Agent), otlptracegrpc.WithInsecure(),
		)
	}

	if err != nil {
		log.Fatalf("failed to initialize export pipeline: %v", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNamespaceKey.String("1995parham"),
			semconv.ServiceNameKey.String("fandogh"),
		),
	)
	if err != nil {
		log.Fatalf("failed to merge resources: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("1995parham.me/fandogh")

	return tracer
}

func Provide(cfg telemetryConfig.Trace) trace.Tracer {
	return New(cfg)
}
