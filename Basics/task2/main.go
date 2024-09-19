package main

// Using Select with Go-Routines and Concurrency

import (
	"fmt"
	"time"
)

func CountUp(c chan int) {
	for i := 0; i <= 10; i++ {
		c <- i
		time.Sleep(time.Microsecond * 500)
	}
	close(c)
}

func CountDown(c chan int) {
	for i := 10; i >= 0; i-- {
		c <- i
	}
	close(c)
}

func main() {
	countUp := make(chan int)
	countDown := make(chan int)

	go CountUp(countUp)
	go CountDown(countDown)

	for {
		select {
		case cU, ok := <-countDown:
			if ok {
				fmt.Println(cU)
			} else {
				countDown = nil
			}
		case cD, ok := <-countUp:
			if ok {
				fmt.Println(cD)
			} else {
				countUp = nil
			}
		}
		if countDown == nil && countUp == nil {
			fmt.Println("Finished Counting")
			break
		}
	}
}
