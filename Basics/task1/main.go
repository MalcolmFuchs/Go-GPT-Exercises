package main

// Using Select with Go-Routines and Concurrency

import (
	"fmt"
	"sync"
	"time"
)

func CountUp(c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(time.Millisecond * 500)
	}

	close(c)
}

func CountDown(c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 10; i > 0; i-- {
		c <- i
		time.Sleep(time.Millisecond * 500)
	}

	close(c)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(2)

	go CountDown(ch, &wg)
	go CountUp(ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for cD := range ch {
		fmt.Println(cD)
	}
}
