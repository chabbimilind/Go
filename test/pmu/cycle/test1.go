// run

package main

import (
    "fmt"
    "log"
    "os"
    "runtime/pprof"
    "sync"
)

//go:noinline
func fun1(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
   fmt.Println(sum)
}

//go:noinline
func fun2(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun3(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun4(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun5(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun6(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun7(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun8(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun9(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

//go:noinline
func fun10(wg *sync.WaitGroup) {
    defer wg.Done()

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

func main() {
    cpuf, err := os.Create("test1_profile")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(cpuf)
    defer pprof.StopCPUProfile()
    
    var wg sync.WaitGroup
    wg.Add(10) // fun1-fun10
    defer wg.Wait() // similar to pthread_join
    
    go fun1(&wg)
    go fun2(&wg)
    go fun3(&wg)
    go fun4(&wg)
    go fun5(&wg)
    go fun6(&wg)
    go fun7(&wg)
    go fun8(&wg)
    go fun9(&wg)
    go fun10(&wg)
}
