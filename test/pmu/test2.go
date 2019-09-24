// run
// Example of usage: go run test2.go

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	_ "time"
)

var wg sync.WaitGroup

var racy int32

func f1() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f2() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f3() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f4() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f5() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f6() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f7() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f8() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f9() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func f10() {
	defer wg.Done()

	var sum int
	for i := 0; i < 500000000; i++ {
		sum -= i / 2
		sum *= i
		sum /= i/3 + 1
		sum -= i / 4
	}

	fmt.Println(sum)
}

func run() error {
	cycleFile, err := os.Create("cycle_profile")
	if err != nil {
		return err
	}
	defer cycleFile.Close()

	var cycle pprof.PMUEventConfig
	cycle.Period = 100000000

	instrFile, err := os.Create("instr_profile")
	if err != nil {
		return err
	}
	defer instrFile.Close()

	var instr pprof.PMUEventConfig
	instr.Period = 100000000

	cacheMissFile, err := os.Create("cacheMiss_profile")
	if err != nil {
		return err
	}
	defer cacheMissFile.Close()

	var cacheMiss pprof.PMUEventConfig
	cacheMiss.Period = 10

	cacheRefFile, err := os.Create("cacheRef_profile")
	if err != nil {
		return err
	}
	defer cacheRefFile.Close()

	var cacheRef pprof.PMUEventConfig
	cacheRef.Period = 1000

	if err := pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(cycleFile, &cycle), pprof.WithProfilingPMUInstructions(instrFile, &instr), pprof.WithProfilingPMUCacheReferences(cacheRefFile, &cacheRef), pprof.WithProfilingPMUCacheMisses(cacheMissFile, &cacheMiss)); err != nil {
		return err
	}

	defer pprof.StopPMUProfile()

	wg.Add(10)
	defer wg.Wait()

	go f1()
	go f2()
	go f3()
	go f4()
	go f5()
	go f6()
	go f7()
	go f8()
	go f9()
	go f10()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
