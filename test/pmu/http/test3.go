// run
// Example of usage: 
// 1. go run test3.go
// 2. go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6\&pmu=true\&pmuevent=cacheMisses\&pmuperiod=10000000

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
	"net/http"
	_ "net/http/pprof"
)

var wg sync.WaitGroup
var mux sync.Mutex
var sum int

func f(i int) {
	defer wg.Done()
	for j := i; j < 100000000; j++ {
        sum -= j / 2
        sum *= j
        time.Sleep(time.Microsecond)
    }
}

func run() error {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	wg.Add(1000)
	defer wg.Wait()

	for i := 0; i < 1000; i++ {
		go f(i)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sum)
}
