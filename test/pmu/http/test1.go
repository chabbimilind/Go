// run
// Example of usage:
// 1. go run test1.go
// 2. go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6\&pmu=true\&pmuevent=cycles\&pmuperiod=10000000

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	_ "time"
)

var sum int

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// time.Sleep(2 * time.Second)
	for i := 0; i <= 10000000000; i++ {
		sum += i
	}
	fmt.Println(sum)
}
