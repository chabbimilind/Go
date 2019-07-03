// run

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
	itimer, err := os.Create("itimer_profile")
	if err != nil {
	 return err
	}
    defer itimer.Close()
    
    cycle, err := os.Create("cycle_profile")
	if err != nil {
	 return err
	}
    defer cycle.Close()
    
    instr, err := os.Create("instr_profile")
	if err != nil {
	 return err
	}
    defer instr.Close()
    
    cacheRef, err := os.Create("cacheRef_profile")
	if err != nil {
	 return err
	}
    defer cacheRef.Close()
	
    cacheMiss, err := os.Create("cacheMiss_profile")
	if err != nil {
	 return err
	}
    defer cacheMiss.Close()
	
    if err := pprof.StartCPUProfile(itimer); err != nil {
        return err
	}
	defer pprof.StopCPUProfile()

    if err := pprof.StartPMUProfile(pprof.WithProfilingCycle(cycle, 20000000), pprof.WithProfilingInstr(instr, 20000000), pprof.WithProfilingCacheRef(cacheRef, 1000), pprof.WithProfilingCacheMiss(cacheMiss, 1)); err != nil {
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
