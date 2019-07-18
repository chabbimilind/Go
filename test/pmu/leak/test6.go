// run

package main

import (
	"fmt"
	"os"
	"time"
	"runtime/pprof"
)

var sum int

func run() {
	for {
		cycleFile, err := os.Create("cycle_profile")
		if err != nil {
			return
		}

		instrFile, err := os.Create("instr_profile")
		if err != nil {
			return
		}

		cacheRefFile, err := os.Create("cacheRef_profile")
		if err != nil {
			return
		}

		cacheMissFile, err := os.Create("cacheMiss_profile")
		if err != nil {
			return
		}

		var cycle pprof.PMUEventConfig
		cycle.Period = 1000000
		var instr pprof.PMUEventConfig
		instr.Period = 1000000
		var cacheRef pprof.PMUEventConfig
		cacheRef.Period = 100
		var cacheMiss pprof.PMUEventConfig
		cacheMiss.Period = 1

		if err := pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(cycleFile, &cycle), pprof.WithProfilingPMUInstructions(instrFile, &instr), pprof.WithProfilingPMUCacheReferences(cacheRefFile, &cacheRef), pprof.WithProfilingPMUCacheMisses(cacheMissFile, &cacheMiss)); err != nil {
			return
		}

		for i := 0; i < 100000000; i++ {
			sum += i
		}
		pprof.StopPMUProfile()

		cycleFile.Close()
		instrFile.Close()
		cacheRefFile.Close()
		cacheMissFile.Close()
	}
}

func main() {
	for i := 0; i < 100; i++ {
		go run()
	}

	time.Sleep(time.Hour)
	fmt.Println(sum)
}
