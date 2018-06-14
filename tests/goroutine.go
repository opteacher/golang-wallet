package main

import (
	"sync"
	"log"
	"fmt"
	"time"
)

var sig = make(chan int)
var end = make(chan bool)
func testChannel() {
	for i := range sig {
		time.Sleep(2000 * time.Millisecond)
		log.Println(i)
	}
	end <- true
	log.Println("End")
}
func testChannelByOK() {
	for {
		i, ok := <- sig
		if !ok { break }
		time.Sleep(2000 * time.Millisecond)
		log.Println(i)
	}
	log.Println("End")
}

func main() {
	log.SetFlags(log.Lshortfile)

	// Test goroutine
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(id int) {
			log.Println(id);
			defer wg.Done()
		}(i)
	}
	wg.Wait()

	// Test channel
	fmt.Println()
	go testChannel()
	for i := 0; i < 5; i++ {
		sig <- i
	}
	close(sig)
	<- end
	close(end)

	//Use of-idiom test channel
	fmt.Println()
	sig = make(chan int)
	go testChannelByOK()
	sig <- 20
	sig <- 30
	sig <- 40
	close(sig)
	time.Sleep(5 * time.Second)
}
