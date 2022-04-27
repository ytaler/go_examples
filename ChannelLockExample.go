package main

import (
	"fmt"
	"time"
)

func main() {
	jobChan := make(chan int, 10)
	go receiveChanInfo(jobChan)
	for i := 0; i < 20; i++ {
		fmt.Printf("[FOR] loop: %d - len channel: %d\n", i, len(jobChan))
		jobChan <- i
		fmt.Printf("[FOR] after channel: %d - len channel: %d\n", i, len(jobChan))
	}
}

func receiveChanInfo(jobChan chan int) {
	// some random startup time
	time.Sleep(3 * time.Second)
	for {
		chanSize := len(jobChan)
		if chanSize <= 0 {
			break
		}
		fmt.Printf("[C] data: %d - len channel: %d\n", <-jobChan, chanSize)
		time.Sleep(time.Second)
	}

}
