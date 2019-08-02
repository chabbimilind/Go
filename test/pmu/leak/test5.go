// run

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

var wg sync.WaitGroup
var sum int

func run() {
	defer wg.Done()
	for i := 0; i < 10000000; i++ {
		sum += i
	}
}

func main() {
	wg.Add(1000)

	cycleFile, err := os.Create("cycle_profile")
	if err != nil {
		log.Fatal(err)
		return
	}
	instrFile, err := os.Create("instr_profile")
	if err != nil {
		log.Fatal(err)
		return
	}
	cacheRefFile, err := os.Create("cacheRef_profile")
	if err != nil {
		log.Fatal(err)
		return
	}
	cacheMissFile, err := os.Create("cacheMiss_profile")
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		go run()
	}

	wg.Wait()
	fmt.Println(sum)
	pprof.StopPMUProfile()

	cycleFile.Close()
	instrFile.Close()
	cacheRefFile.Close()
	cacheMissFile.Close()
}
