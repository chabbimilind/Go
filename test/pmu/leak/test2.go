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
	for /*i := 0; i < 100000; i++*/ {
		cycleFile, err := os.Create("cycle_profile")
		if err != nil {
			return err
		}

		var cycle pprof.PMUEventConfig
		cycle.Period = 100000
		cycle.IsKernelIncluded = false
		cycle.IsHvIncluded = false

		if err = pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(cycleFile, &cycle)); err != nil {
			return err
		}
		for j := 0; j < 10000000; j++ {
			sum += j
		}
		pprof.StopPMUProfile()
		cycleFile.Close()
	}
	fmt.Println(sum)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
