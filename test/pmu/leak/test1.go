// run

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var sum int

func run() error {
	for /*i := 0; i < 10; i++*/ {
		itimerFile, err := os.Create("itimer_profile")
		if err != nil {
			return err
		}

		if err = pprof.StartCPUProfile(itimerFile); err != nil {
			return err
		}
		for j := 0; j < 100; j++ {
			sum += j
		}
		pprof.StopCPUProfile()
		itimerFile.Close()
	}
	fmt.Println(sum)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
