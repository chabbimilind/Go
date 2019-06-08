// run

package main

import (
    "fmt"
    "log"
    "os"
    "runtime/pprof"
)

//go:noinline
func fun1(c chan int) {
    defer close(c)

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
func fun2(c chan int) {
    defer close(c)

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
func fun3(c chan int) {
    defer close(c)

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
func fun4(c chan int) {
    defer close(c)

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
func fun5(c chan int) {
    defer close(c)

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
func fun6(c chan int) {
    defer close(c)

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
func fun7(c chan int) {
    defer close(c)

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
func fun8(c chan int) {
    defer close(c)

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
func fun9(c chan int) {
    defer close(c)

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
func fun10(c chan int) {
    defer close(c)

    var sum int
    for i:= 0; i < 500000000; i++ {
        sum -= i / 2
        sum *= i
        sum /= i / 3 + 1
        sum -= i / 4
    }
    
    fmt.Println(sum)
}

func waitForChans(chans ...chan int) {
    for _, v := range chans {
		<-v
	}
}

func main() {
    cpuf, err := os.Create("test2_profile")
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(cpuf)
    defer pprof.StopCPUProfile()
    
    var chans [10]chan int
    for i := range chans {
        chans[i] = make(chan int)
    }
    
    go fun1(chans[0])
    go fun2(chans[1])
    go fun3(chans[2])
    go fun4(chans[3])
    go fun5(chans[4])
    go fun6(chans[5])
    go fun7(chans[6])
    go fun8(chans[7])
    go fun9(chans[8])
    go fun10(chans[9])

    waitForChans(chans[0], chans[1], chans[2], chans[3], chans[4], chans[5], chans[6], chans[7], chans[8], chans[9])
}
