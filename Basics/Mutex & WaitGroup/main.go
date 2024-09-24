// Mutex & WaitGroup
package main

import (
	"fmt"
	"sync"
)

var (
	mutex sync.Mutex
	count int
)

func increment() {
	mutex.Lock()
	defer mutex.Unlock()
	count++
	fmt.Printf("Zählerstand nach Inkrement: %d\n", count)
}

func main() {
	var wg sync.WaitGroup
	maxAmount := 5

	for i := 0; i < maxAmount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Finaler Zählerstand: %d\n", count)
}
