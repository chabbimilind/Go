// run
// Example of Usage: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=6\&pmu=true\&event=r53010e\&period=1000000

package main

import (
	"fmt"
	"log"
	// "time"
	"net/http"
	_ "net/http/pprof"
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
