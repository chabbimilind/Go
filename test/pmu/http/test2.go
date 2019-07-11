// run
// Example of Usage: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6\&pmu=true\&event=cycles\&period=10000000

package main

import (
	"fmt"
	"log"
	"sync"
    // "time"
    "net/http"
    _ "net/http/pprof"
)

var wg sync.WaitGroup
var mux sync.Mutex
var sum int

func f(i int) {
	defer wg.Done()
        var local int
	for j := i; j < 100000000; j++ {
        local -= j / 2
        local *= j
        mux.Lock()
        sum += local
        mux.Unlock()
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
