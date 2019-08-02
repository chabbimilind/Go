// run
// Example of Usage:
// 1. go run test4.go
// 2. go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6\&pmu=true\&pmuevent=r53010e\&pmuperiod=1000000

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var sum int

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	for i := 0; i <= 10000000000; i++ {
		sum += i
	}
	fmt.Println(sum)
}
