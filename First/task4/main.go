package main

//time.After und ctx

import (
	"context"
	"fmt"
	"time"
)

func CountUp(c chan int, ctx context.Context) {
	defer close(c)

	for i := 0; i <= 10; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("CountDown: Abbruch Signal erhalten, Beende Goroutine")
			return
		default:
			c <- i
			time.Sleep(time.Millisecond * 500)
		}
	}

}

func CountDown(c chan int, ctx context.Context) {
	defer close(c)

	for i := 10; i >= 0; i-- {
		select {
		case <-ctx.Done():
			fmt.Println("CountDown: Abbruch Signal erhalten, Beende Goroutine")
			return
		default:
			c <- i
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2000)
	defer cancel()

	countUp := make(chan int)
	countDown := make(chan int)

	go CountUp(countUp, ctx)
	go CountDown(countDown, ctx)

	for {
		select {
		case cD, ok := <-countDown:
			if ok {
				fmt.Println(cD)
			} else {
				countDown = nil
			}
		case cU, ok := <-countUp:
			if ok {
				fmt.Println(cU)
			} else {
				countUp = nil
			}
		case <-ctx.Done():
			fmt.Println("Timeout: Keine weiteren Daten empfangen.")
			return
		}
		if countDown == nil && countUp == nil {
			fmt.Println("Finished Counting!")
			break
		}
	}

}
