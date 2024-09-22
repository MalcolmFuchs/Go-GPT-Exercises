// Metrics und Distributed Tracing mit OpenTelemetry
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() (func(), error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("web-crawler"),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background()) }, nil
}

func crawlWebsite(ctx context.Context, url string, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "crawl-website")
	span.SetAttributes(attribute.String("website.url", url))
	defer span.End()

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("crawl.status", "error"))
		log.Printf("Fehler beim Crawlen von %s: %vn\n", url, err)
	} else {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		span.SetAttributes(attribute.String())
	}
}

func main() {

}
