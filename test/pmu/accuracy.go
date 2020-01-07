/* run
Example of usage: go build accuracy.go && ./accuracy -p=false -r=10000 -m=10000 && go tool pprof -top my_prof && ./accuracy -p=true -i=10000 -m=10000 && go tool pprof -top my_prof
Example output showing the difference in the profiling accuracy:
---------- begin output ------------
doPMU= false
multiplier= 10000
rate= 10000
interval= 1000000
1029153122977
File: accuracy
Type: cpu
Time: Nov 12, 2019 at 5:48am (UTC)
Duration: 200.46ms, Total samples = 1300us ( 0.65%)
Showing nodes accounting for 1300us, 100% of 1300us total
      flat  flat%   sum%        cum   cum%
     300us 23.08% 23.08%      300us 23.08%  main.C_expect_5_46
     300us 23.08% 46.15%      300us 23.08%  main.G_expect_12_73
     200us 15.38% 61.54%      200us 15.38%  main.F_expect_10_91
     200us 15.38% 76.92%      200us 15.38%  main.H_expect_14_546
     200us 15.38% 92.31%      200us 15.38%  main.J_expect_18_18
     100us  7.69%   100%      100us  7.69%  main.B_expect_3_64
         0     0%   100%     1300us   100%  main.main
         0     0%   100%     1300us   100%  runtime.main
doPMU= true
multiplier= 10000
rate= 100
interval= 10000
1029153122977
File: accuracy
Type: cycles
Time: Nov 12, 2019 at 5:48am (UTC)
Showing nodes accounting for 137240000, 99.10% of 138490000 total
Dropped 107 nodes (cum <= 692450)
      flat  flat%   sum%        cum   cum%
  24690000 17.83% 17.83%   24700000 17.84%  main.J_expect_18_18
  22340000 16.13% 33.96%   22350000 16.14%  main.I_expect_16_36
  19810000 14.30% 48.26%   19820000 14.31%  main.H_expect_14_546
  17240000 12.45% 60.71%   17240000 12.45%  main.G_expect_12_73
  14620000 10.56% 71.27%   14620000 10.56%  main.F_expect_10_91
  12280000  8.87% 80.14%   12280000  8.87%  main.E_expect_9_09
   9930000  7.17% 87.31%    9930000  7.17%  main.D_expect_7_27
   7320000  5.29% 92.59%    7320000  5.29%  main.C_expect_5_46
   4900000  3.54% 96.13%    4900000  3.54%  main.B_expect_3_64
   2440000  1.76% 97.89%    2450000  1.77%  main.A_expect_1_82
    660000  0.48% 98.37%    1790000  1.29%  runtime/pprof.(*profMap).lookup
    590000  0.43% 98.79%     960000  0.69%  runtime.mapaccess1_fast64
    340000  0.25% 99.04%    2130000  1.54%  runtime/pprof.(*profileBuilder).addData
     80000 0.058% 99.10%  135700000 97.99%  main.main
         0     0% 99.10%  135700000 97.99%  runtime.main
         0     0% 99.10%    2130000  1.54%  runtime/pprof.(*profileBuilder).addPMUData
         0     0% 99.10%    2600000  1.88%  runtime/pprof.pmuProfileWriter
---------- end output ------------
*/
package main

import (
	"fmt"
	//"math/rand"
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

func J_expect_18_18(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func I_expect_16_36(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func H_expect_14_546(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func G_expect_12_73(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func F_expect_10_91(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func E_expect_9_09(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func D_expect_7_27(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func C_expect_5_46(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func B_expect_3_64(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func A_expect_1_82(v uint64, trip uint64) uint64 {
	ret := v
	for i := trip; i > 0; i-- {
		ret += i
		ret = ret ^ (i + 0xcafebabe)
	}
	return ret
}
func main() {
	var q uint64
	// PMU workd, wall fails multiplier := uint64(1000000)
	//multiplier := uint64(1000000)
	//multiplier := uint64(10000)
	// Should we profile manually?
	doPMU := flag.Bool("p", false, "do PMU profiling")
	rate := flag.Int("r", 100, "Wallclock profile sampling rate")
	interval := flag.Int64("i", 1000000, "PMU event interval")
	multiplier := flag.Uint64("m", 1000000, "multiplier")
	flag.Parse()
	fmt.Println("doPMU=", *doPMU)
	fmt.Println("multiplier=", *multiplier)
	fmt.Println("rate=", *rate)
	fmt.Println("interval=", *interval)

	file, err := os.Create("my_prof")
	if err != nil {
		log.Fatal(err)
	}
	if *doPMU {
		var cycle pprof.PMUEventConfig
		cycle.Period = *interval
		cycle.IsKernelIncluded = false
		cycle.IsHvIncluded = false

		if err = pprof.StartPMUProfile(pprof.WithProfilingPMUCycles(file, &cycle)); err != nil {
			log.Fatal(err)
		}
	} else {
		if err = pprof.StartCPUProfile(file, *rate); err != nil {
			log.Fatal(err)
		}
	}
	mult := *multiplier

	for i := uint64(0); i < 100; i++ {
		f := i + A_expect_1_82(0xebabefac23, 1*mult)
		//fmt.Println(f)
		g := i + B_expect_3_64(f, 2*mult)
		//fmt.Println(g)
		h := i + C_expect_5_46(g, 3*mult)
		//fmt.Println(h)
		k := i + D_expect_7_27(h, 4*mult)
		//fmt.Println(k)
		l := i + E_expect_9_09(k, 5*mult)
		//fmt.Println(l)
		m := i + F_expect_10_91(l, 6*mult)
		//fmt.Println(m)
		n := i + G_expect_12_73(m, 7*mult)
		//fmt.Println(n)
		o := i + H_expect_14_546(n, 8*mult)
		//fmt.Println(o)
		p := i + I_expect_16_36(o, 9*mult)
		//fmt.Println(p)
		q = i + J_expect_18_18(p, 10*mult)
		//fmt.Println(q)
	}
	if *doPMU {
		pprof.StopPMUProfile()
	} else {
		pprof.StopCPUProfile()
	}
	file.Close()
	fmt.Println(q)

}
