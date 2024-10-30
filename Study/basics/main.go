package main

import (
	"fmt"
	"sync"
)

type Worker struct {
	ID int
}

type WorkerInterface interface {
	DoWork(ch chan int)
}

func (w Worker) DoWork(ch chan int) {
	ch <- w.ID
}

func main() {
	count := 10
	ch := make(chan int, count)
	var wg sync.WaitGroup

	workers := make([]WorkerInterface, count)
	for i := 0; i < count; i++ {
		workers[i] = Worker{ID: i}
	}

	for _, worker := range workers {
		wg.Add(1)

		go func(w WorkerInterface) {
			defer wg.Done()
			w.DoWork(ch)
		}(worker)
	}

	wg.Wait()
	close(ch)

	for num := range ch {
		fmt.Println(num)
	}

}
