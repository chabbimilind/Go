// run
// Example of usage: go run test1.go

package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"runtime/pprof"
)

var wg sync.WaitGroup

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
	cycle.Period =  10000000
	cycle.IsKernelIncluded = false
	cycle.IsHvIncluded = false

	if err := pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(cycleFile, &cycle)); err != nil {
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
