package main

// time.After

import (
	"fmt"
	"time"
)

func CountUp(c chan int) {
	for i := 0; i <= 10; i++ {
		c <- i
		time.Sleep(time.Millisecond * 500)
	}
	close(c)
}

func CountDown(c chan int) {
	for i := 10; i >= 0; i-- {
		c <- i
		time.Sleep(time.Millisecond * 500)
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
		case cU, ok := <-countUp:
			if ok {
				fmt.Println(cU)
			} else {
				countUp = nil
			}
		case cD, ok := <-countDown:
			if ok {
				fmt.Println(cD)
			} else {
				countDown = nil
			}
		case <-time.After(time.Millisecond * 1000):
			fmt.Println("Timeout: Keine Daten innerhalb von 1 Sekunde empfangen.")
		}
		if countDown == nil && countUp == nil {
			fmt.Println("Finished Counting")
			break
		}
	}

}
