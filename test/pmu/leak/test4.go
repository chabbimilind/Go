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
	for /*i := 0; i < 100; i++*/ {
		cycleFile, err := os.Create("cycle_profile")
		if err != nil {
			return err
		}
		instrFile, err := os.Create("instr_profile")
		if err != nil {
			return err
		}
		cacheRefFile, err := os.Create("cacheRef_profile")
		if err != nil {
			return err
		}
		cacheMissFile, err := os.Create("cacheMiss_profile")
		if err != nil {
			return err
		}

		var cycle pprof.PMUEventConfig
		cycle.Period = 1000000
		var instr pprof.PMUEventConfig
		instr.Period = 1000000
		var cacheRef pprof.PMUEventConfig
		cacheRef.Period = 10000
		var cacheMiss pprof.PMUEventConfig
		cacheMiss.Period = 1000

		if err := pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(cycleFile, &cycle), pprof.WithProfilingPMUInstructions(instrFile, &instr), pprof.WithProfilingPMUCacheReferences(cacheRefFile, &cacheRef), pprof.WithProfilingPMUCacheMisses(cacheMissFile, &cacheMiss)); err != nil {
			return err
		}
		for j := 0; j < 10000000; j++ {
			sum += j
		}
		pprof.StopPMUProfile()

		cycleFile.Close()
		instrFile.Close()
		cacheRefFile.Close()
		cacheMissFile.Close()
	}
	fmt.Println(sum)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
