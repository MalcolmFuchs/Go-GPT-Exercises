// Logging & Observability mit zusätzlichen Metriken

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

func RetryMechanismen(web Website, maxRetries int, retryDur time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	for i := 0; i < maxRetries; i++ {
		log.Printf("Retry %d für %s\n", i+1, web.URL)

		err := web.Crawl(ctx)
		if err == nil {
			return nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout beim Crawlen von %s nach Retry %d\n", web.URL, i+1)
			atomic.AddUint64(&timeoutCount, 1)
		} else {
			log.Printf("Fehler beim Crawlen von %s nach Retry %d: %s\n", web.URL, i+1, err)
			atomic.AddUint64(&errorCount, 1)
		}

		if i == maxRetries-1 {
			return err
		}

		time.Sleep(retryDur)
	}

	return errors.New("Alle Crawl-Versuche fehlgeschlagen")
}

func handleCrawl(web Website, sem chan struct{}, wg *sync.WaitGroup, errChan chan error) {
	start := time.Now()
	defer wg.Done()

	sem <- struct{}{}

	err := RetryMechanismen(web, 3, time.Second)
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
	defer logMemoryUsage()

	var wg sync.WaitGroup
	maxCrawlers := 3

	Websites := []Website{
		{URL: "https://example1.com"},
		{URL: "https://example2.com"},
		{URL: "https://example3.com"},
		{URL: "https://example4.com"},
		{URL: "https://example5.com"},
	}

	sem := make(chan struct{}, maxCrawlers)
	errChan := make(chan error, len(Websites))

	for _, web := range Websites {
		wg.Add(1)
		go handleCrawl(web, sem, &wg, errChan)
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
