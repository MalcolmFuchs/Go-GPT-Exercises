// Not Working yet...
// Metrics und Distributed Tracing mit OpenTelemetry (angepasst für Prometheus und HTTP Handler)
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
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	meter            metric.Meter
	successCounter   metric.Int64Counter
	errorCounter     metric.Int64Counter
	latencyHistogram metric.Float64Histogram
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

func initMetrics() (*prometheus.Exporter, error) {
	exporter, err := prometheus.New(prometheus.Config{})
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(meterProvider)

	meter = meterProvider.Meter("web-crawler")

	// Instrumente registrieren
	successCounter, _ = meter.Int64Counter(
		"crawler.success_count",
		metric.WithDescription("Anzahl der erfolgreichen Crawls"),
	)

	errorCounter, _ = meter.Int64Counter(
		"crawler.error_count",
		metric.WithDescription("Anzahl der fehlgeschlagenen Crawls"),
	)

	latencyHistogram, _ = meter.Float64Histogram(
		"crawler.latency",
		metric.WithDescription("Latenzzeit der Crawls"),
	)

	return exporter, nil
}

func crawlWebsite(ctx context.Context, url string, tracer trace.Tracer) {
	ctx, span := tracer.Start(ctx, "crawl-website")
	span.SetAttributes(attribute.String("website.url", url))
	defer span.End()

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	// Metrik für Latenz erfassen
	latencyHistogram.Record(ctx, duration.Seconds())

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("crawl.status", "error"))
		errorCounter.Add(ctx, 1) // Fehler-Counter erhöhen
		log.Printf("Fehler beim Crawlen von %s: %v\n", url, err)
	} else {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		successCounter.Add(ctx, 1) // Erfolgs-Counter erhöhen
		log.Printf("Erfolgreich gecrawlt: %s\n", url)
	}
}

func main() {
	cleanupTracer, err := initTracer()
	if err != nil {
		log.Fatalf("Fehler beim Initialisieren des Tracers: %v", err)
	}
	defer cleanupTracer()

	promExporter, err := initMetrics()
	if err != nil {
		log.Fatalf("Fehler beim Initialisieren der Metriken: %v", err)
	}

	urls := []string{
		"https://example.com",
		"https://example2.com",
		"https://example3.com",
		"https://nonexistent.url",
	}

	tracer := otel.Tracer("web-crawler")

	for _, url := range urls {
		crawlWebsite(context.Background(), url, tracer)
	}

	// Starte einen HTTP-Server, um Prometheus-Metriken verfügbar zu machen
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promExporter.ServeHTTP(w, r)
	})
	log.Println("Prometheus Metriken verfügbar unter :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
