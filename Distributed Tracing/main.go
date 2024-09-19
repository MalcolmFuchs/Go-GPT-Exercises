// Distributed Tracing

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type Website struct {
	URL string
}

var errorCount uint64
var successCount uint64
var timeoutCount uint64
var totalDuration uint64

func (web Website) Crawl(ctx context.Context) error {
	log.Printf("Starte das Crawlen der Website %s\n", web.URL)

	select {
	case <-time.After(time.Millisecond * 500):
		if web.URL == "https://example4.com" {
			return errors.New("Fehler beim Crawlen der Webseite: " + web.URL)
		}
		log.Printf("Das Crawling der Webseite %s wurde erfolgreich abgeschlossen\n", web.URL)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("Timeout erreicht: %w", ctx.Err())
	}
}

func RetryMechanismen(ctx context.Context, web Website, maxRetries int, retryDur time.Duration) error {

	tracer := otel.Tracer("Crawl-Operation")

	for i := 0; i < maxRetries; i++ {

		ctx, span := tracer.Start(ctx, fmt.Sprintf("Retry %d für %s", i+1, web.URL))
		defer span.End()

		log.Printf("Retry %d für %s\n", i+1, web.URL)

		err := web.Crawl(ctx)
		if err == nil {
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(200))
			return nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout beim Crawlen von %s nach Retry %d\n", web.URL, i+1)
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(408))
			atomic.AddUint64(&timeoutCount, 1)
		} else {
			log.Printf("Fehler beim Crawlen von %s nach Retry %d: %s\n", web.URL, i+1, err)
			span.SetAttributes(semconv.HTTPStatusCodeKey.Int(500))
			atomic.AddUint64(&errorCount, 1)
		}

		if i == maxRetries-1 {
			return err
		}

		time.Sleep(retryDur)
	}

	return errors.New("Alle Crawl-Versuche fehlgeschlagen")
}

func handleCrawl(ctx context.Context, web Website, sem chan struct{}, wg *sync.WaitGroup, errChan chan error) {
	start := time.Now()
	defer wg.Done()

	sem <- struct{}{}

	err := RetryMechanismen(ctx, web, 3, time.Second)
	if err != nil {
		errChan <- err
	} else {
		atomic.AddUint64(&successCount, 1)
		errChan <- nil
	}

	duration := time.Since(start)
	atomic.AddUint64(&totalDuration, uint64(duration.Nanoseconds()))
	log.Printf("Dauer der Anfrage: %v\n", duration)

	<-sem
}

func main() {
	Websites := []Website{
		{URL: "https://example1.com"},
		{URL: "https://example2.com"},
		{URL: "https://example3.com"},
		{URL: "https://example4.com"},
		{URL: "https://example5.com"},
	}

	var wg sync.WaitGroup
	maxCrawlers := 3
	sem := make(chan struct{}, maxCrawlers)
	errChan := make(chan error, len(Websites))
	cleanup := initTracer()

	defer logMemoryUsage()
	defer cleanup()

	rootCtx := context.Background()

	for _, web := range Websites {
		wg.Add(1)
		go handleCrawl(rootCtx, web, sem, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			log.Println("Fehler: ", err)
		}
	}

	log.Printf("Anzahl der Fehler: %d\n", errorCount)
	log.Printf("Anzahl der Erfolgreichen: %d\n", successCount)
	log.Printf("Anzahl der Timeouts: %d\n", timeoutCount)
	log.Printf("Durchschnittliche Anfragedauer: %v\n", time.Duration(totalDuration)/time.Duration(successCount+errorCount))
	log.Println("Alle Webseiten wurden verarbeitet.")
}

func logMemoryUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Printf("Speicherverbrauch: %v KB", memStats.Alloc/1024)
	log.Printf("Heap-Alloc: %v KB", memStats.HeapAlloc/1024)
	log.Printf("GC-Zyklen: %v", memStats.NumGC)
}

func initTracer() func() {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Crawl-Service"),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}
