// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CPU profiling.
//
// The signal handler for the profiling clock tick adds a new stack trace
// to a log of recent traces. The log is read by a user goroutine that
// turns it into formatted profile data. If the reader does not keep up
// with the log, those writes will be recorded as a count of lost records.
// The actual profile buffer is in profbuf.go.

package runtime

import (
	"runtime/internal/atomic"
	"runtime/internal/sys"
	"unsafe"
)

const maxCPUProfStack = 64
// MaxPMUEvent is a small number and cannot be >= 10 because we do a linear search on it
const MaxPMUEvent = 10

type profile struct {
	lock mutex
	on   bool     // profiling is on
	log  *profBuf // profile events written here

	// extra holds extra stacks accumulated in addNonGo
	// corresponding to profiling signals arriving on
	// non-Go-created threads. Those stacks are written
	// to log the next time a normal Go thread gets the
	// signal handler.
	// Assuming the stacks are 2 words each (we don't get
	// a full traceback from those threads), plus one word
	// size for framing, 100 Hz profiling would generate
	// 300 words per second.
	// Hopefully a normal Go thread will get the profiling
	// signal at least once every few seconds.
	extra     [1000]uintptr
	numExtra  int
	lostExtra uint64 // count of frames lost because extra is full
}

var cpuprof profile
var pmuprof [MaxPMUEvent]profile

// SetCPUProfileRate sets the CPU profiling rate to hz samples per second.
// If hz <= 0, SetCPUProfileRate turns off profiling.
// If the profiler is on, the rate cannot be changed without first turning it off.
//
// Most clients should use the runtime/pprof package or
// the testing package's -test.cpuprofile flag instead of calling
// SetCPUProfileRate directly.
func SetCPUProfileRate(hz int) {
	// Clamp hz to something reasonable.
	if hz < 0 {
		hz = 0
	}
	if hz > 1000000 {
		hz = 1000000
	}

	lock(&cpuprof.lock)
	defer unlock(&cpuprof.lock)
	if hz > 0 {
		if cpuprof.on || cpuprof.log != nil {
			print("runtime: cannot set cpu profile rate until previous profile has finished.\n")
			return
		}

		cpuprof.on = true
		// Enlarging the buffer words and tags reduces the number of samples lost at the cost of larger amounts of memory 
		cpuprof.log = newProfBuf(/* header size */ 1, /* buffer words */ 1<<17, /* tags */ 1<<14)
		hdr := [1]uint64{uint64(hz)}
		cpuprof.log.write(nil, nanotime(), hdr[:], nil)
		setcpuprofilerate(int32(hz))
	} else if cpuprof.on {
		setcpuprofilerate(0)
		cpuprof.on = false
		cpuprof.addExtra()
		cpuprof.log.close()
	}
}

func SetPMUProfile(eventId int, eventAttr *PMUEventAttr) {
	lock(&pmuprof[eventId].lock)
	defer unlock(&pmuprof[eventId].lock)
	if eventAttr != nil {
		if pmuprof[eventId].on || pmuprof[eventId].log != nil {
			print("runtime: cannot set pmu profile rate until previous profile has finished.\n")
			return
		}

		pmuprof[eventId].on = true
		// Enlarging the buffer words and tags reduces the number of samples lost at the cost of larger amounts of memory 
		pmuprof[eventId].log = newProfBuf(/* header size */ 1, /* buffer words */ 1<<17, /* tags */ 1<<14)
		hdr := [1]uint64{eventAttr.Period}
		pmuprof[eventId].log.write(nil, nanotime(), hdr[:], nil)
		setpmuprofile(int32(eventId), eventAttr)
	} else if pmuprof[eventId].on {
		setpmuprofile(int32(eventId), nil)
		pmuprof[eventId].on = false
		pmuprof[eventId].addExtra(eventId)
		pmuprof[eventId].log.close()
	}
}

//go:nowritebarrierrec
func (p *profile) addImpl(gp *g, stk []uintptr, cpuorpmuprof *profile) {
	if p.numExtra > 0 || p.lostExtra > 0 {
		p.addExtra()
	}
	hdr := [1]uint64{1}
	// Note: write "knows" that the argument is &gp.labels,
	// because otherwise its write barrier behavior may not
	// be correct. See the long comment there before
	// changing the argument here.
	cpuorpmuprof.log.write(&gp.labels, nanotime(), hdr[:], stk)
}

// add adds the stack trace to the profile.
// It is called from signal handlers and other limited environments
// and cannot allocate memory or acquire locks that might be
// held at the time of the signal, nor can it use substantial amounts
// of stack.
//go:nowritebarrierrec
func (p *profile) add(gp *g, stk []uintptr, eventIds ...int) {
	if len(eventIds) == 0 {
		for !atomic.Cas(&prof.signalLock, 0, 1) {
			osyield()
		}
		if prof.hz != 0 { // implies cpuprof.log != nil
			p.addImpl(gp, stk, &cpuprof)
		}
		atomic.Store(&prof.signalLock, 0)
	} else {
		eventId := eventIds[0]
		for !atomic.Cas(&pmuEvent[eventId].signalLock, 0, 1) {
			osyield()
		}
		if pmuEvent[eventId].eventAttr != nil { // implies pmuprof[eventId].log != nil
			p.addImpl(gp, stk, &pmuprof[eventId])
		}
		atomic.Store(&pmuEvent[eventId].signalLock, 0)
	}
}

//go:nosplit
//go:nowritebarrierrec
func (p *profile) addNonGoImpl(stk []uintptr, prof *profile) {
	if prof.numExtra+1+len(stk) < len(prof.extra) {
		i := prof.numExtra
		prof.extra[i] = uintptr(1 + len(stk))
		copy(prof.extra[i+1:], stk)
		prof.numExtra += 1 + len(stk)
	} else {
		prof.lostExtra++
	}
}

// addNonGo adds the non-Go stack trace to the profile.
// It is called from a non-Go thread, so we cannot use much stack at all,
// nor do anything that needs a g or an m.
// In particular, we can't call cpuprof.log.write.
// Instead, we copy the stack into cpuprof.extra,
// which will be drained the next time a Go thread
// gets the signal handling event.
//go:nosplit
//go:nowritebarrierrec
func (p *profile) addNonGo(stk []uintptr, eventIds ...int) {
	if len(eventIds) == 0 {
		// Simple cas-lock to coordinate with SetCPUProfileRate.
		// (Other calls to add or addNonGo should be blocked out
		// by the fact that only one SIGPROF can be handled by the
		// process at a time. If not, this lock will serialize those too.)
		for !atomic.Cas(&prof.signalLock, 0, 1) {
			osyield()
		}
		p.addNonGoImpl(stk, &cpuprof)
		atomic.Store(&prof.signalLock, 0)
	} else {
		eventId := eventIds[0]
		// Only one SIGPROF for each PMU event can be handled by the process at a time.
		for !atomic.Cas(&pmuEvent[eventId].signalLock, 0, 1) {
			osyield()
		}
		p.addNonGoImpl(stk, &pmuprof[eventId])
		atomic.Store(&pmuEvent[eventId].signalLock, 0)
	}
}

// addExtra adds the "extra" profiling events,
// queued by addNonGo, to the profile log.
// addExtra is called either from a signal handler on a Go thread
// or from an ordinary goroutine; either way it can use stack
// and has a g. The world may be stopped, though.
func (p *profile) addExtra(eventIds ...int) {
	// Copy accumulated non-Go profile events.
	hdr := [1]uint64{1}
	for i := 0; i < p.numExtra; {
		p.log.write(nil, 0, hdr[:], p.extra[i+1:i+int(p.extra[i])])
		i += int(p.extra[i])
	}
	p.numExtra = 0

	// Report any lost events.
	if p.lostExtra > 0 {
		hdr := [1]uint64{p.lostExtra}
		lostStk := [2]uintptr{
			funcPC(_LostExternalCode) + sys.PCQuantum,
			funcPC(_ExternalCode) + sys.PCQuantum,
		}
		if len(eventIds) == 0 {
			cpuprof.log.write(nil, 0, hdr[:], lostStk[:])
		} else {
			eventId := eventIds[0]
			pmuprof[eventId].log.write(nil, 0, hdr[:], lostStk[:])
		}
		p.lostExtra = 0
	}
}

func (p *profile) addLostAtomic64(count uint64, eventIds ...int) {
	hdr := [1]uint64{count}
	lostStk := [2]uintptr{
		funcPC(_LostSIGPROFDuringAtomic64) + sys.PCQuantum,
		funcPC(_System) + sys.PCQuantum,
	}
	if len(eventIds) == 0 {
		cpuprof.log.write(nil, 0, hdr[:], lostStk[:])
	} else {
		eventId := eventIds[0]
		pmuprof[eventId].log.write(nil, 0, hdr[:], lostStk[:])
	}
}

// CPUProfile panics.
// It formerly provided raw access to chunks of
// a pprof-format profile generated by the runtime.
// The details of generating that format have changed,
// so this functionality has been removed.
//
// Deprecated: Use the runtime/pprof package,
// or the handlers in the net/http/pprof package,
// or the testing package's -test.cpuprofile flag instead.
func CPUProfile() []byte {
	panic("CPUProfile no longer available")
}

//go:linkname runtime_pprof_runtime_cyclesPerSecond runtime/pprof.runtime_cyclesPerSecond
func runtime_pprof_runtime_cyclesPerSecond() int64 {
	return tickspersecond()
}

func readProfileImpl(prof *profile) ([]uint64, []unsafe.Pointer, bool) {
	lock(&prof.lock)
	log := prof.log
	unlock(&prof.lock)
	data, tags, eof := log.read(profBufBlocking)
	if len(data) == 0 && eof {
		lock(&prof.lock)
		prof.log = nil
		unlock(&prof.lock)
	}
	return data, tags, eof
}

// readProfile, provided to runtime/pprof, returns the next chunk of
// binary CPU profiling stack trace data, blocking until data is available.
// If profiling is turned off and all the profile data accumulated while it was
// on has been returned, readProfile returns eof=true.
// The caller must save the returned data and tags before calling readProfile again.
//
//go:linkname runtime_pprof_readProfile runtime/pprof.readProfile
func runtime_pprof_readProfile(eventIds ...int) ([]uint64, []unsafe.Pointer, bool) {
	if len(eventIds) == 0 {
		return readProfileImpl(&cpuprof)
	} else {
		eventId := eventIds[0]
		return readProfileImpl(&pmuprof[eventId])
	}
}
