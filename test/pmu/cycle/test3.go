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
    cpuf, err := os.Create("test3_profile")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(cpuf)
    defer pprof.StopCPUProfile()
    
    var wg [10]sync.WaitGroup
    for i := range wg {
        wg[i].Add(1)
    }
    
    go fun1(&(wg[0]))
    go fun2(&(wg[1]))
    go fun3(&(wg[2]))
    go fun4(&(wg[3]))
    go fun5(&(wg[4]))
    go fun6(&(wg[5]))
    go fun7(&(wg[6]))
    go fun8(&(wg[7]))
    go fun9(&(wg[8]))
    go fun10(&(wg[9]))
    
    for i:= range wg {
        wg[i].Wait()
    }
}
