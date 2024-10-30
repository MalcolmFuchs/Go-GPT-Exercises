// Retry-Mechanismen & Timeouts

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Website struct {
	URL string
}

func (web Website) Crawl() error {
	fmt.Printf("Starte das Crawlen der Website %s\n", web.URL)
	time.Sleep(time.Millisecond * 500)

	if web.URL == "https://example4.com" {
		return errors.New("Fehler beim beginnes des Crawlens von: " + web.URL)
	}

	fmt.Printf("Das Crawling der Webseite %s wurde erfolgreich ausgeführt \n", web.URL)
	return nil
}

func RetryMechanismen(web Website, maxRetries int, retryDelay time.Duration) error {

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Retry %d für %s\n", i+1, web.URL)

		err := web.Crawl()
		if err == nil {
			return nil
		}

		if i == maxRetries-1 {
			return err
		}

		time.Sleep(retryDelay)

	}

	return errors.New("Alle Crawl-Versuche fehlgeschlagen")
}

func handleCrawl(web Website, sem chan struct{}, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	sem <- struct{}{}

	err := web.Crawl()
	if err != nil {
		retryErr := RetryMechanismen(web, 3, time.Second)
		errChan <- retryErr
	} else {
		errChan <- nil
	}

	<-sem
}

func main() {
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
			fmt.Println("Fehler: ", err)
		}
	}
}
