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

const maxPMUProfStack = 64
const maxPMUEvents  = 10

type pmuProfile struct {
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

var pmuprof [maxPMUEvents]pmuProfile // event -> cpuProfile

func SetPMUProfilePeriod(eventId int, period int) {
    // Clamp period to something reasonable.
    if period < 0 {
        period = 0
    }
    if period > 0 && period < 300 { // follow what hpctoolkit did
        period = 300
    }
    lock(&pmuprof[eventId].lock)
    if period > 0 {
        if pmuprof[eventId].on || pmuprof[eventId].log != nil {
            print("runtime: cannot set pmu profile rate until previous profile has finished.\n")
            unlock(&pmuprof[eventId].lock)
            return
        }

        pmuprof[eventId].on = true
        pmuprof[eventId].log = newProfBuf(1, 1<<17, 1<<14)
        hdr := [1]uint64{uint64(period)}
        pmuprof[eventId].log.write(nil, nanotime(), hdr[:], nil)
        setpmuprofileperiod(int32(eventId), int32(period))
    } else if pmuprof[eventId].on {
        setpmuprofileperiod(int32(eventId), 0)
        pmuprof[eventId].on = false
        pmuprof[eventId].addExtra(eventId)
        pmuprof[eventId].log.close()
    }
    unlock(&pmuprof[eventId].lock)
}

// add adds the stack trace to the profile.
// It is called from signal handlers and other limited environments
// and cannot allocate memory or acquire locks that might be
// held at the time of the signal, nor can it use substantial amounts
// of stack.
//go:nowritebarrierrec
func (p *pmuProfile) add(gp *g, stk []uintptr, eventId int) {
	// Simple cas-lock to coordinate with setpmuprofilerate.
	for !atomic.Cas(&profs[eventId].signalLock, 0, 1) {
		osyield()
	}

	if profs[eventId].period != 0 { // implies pmuprof[eventId].log != nil
		if p.numExtra > 0 || p.lostExtra > 0 {
			p.addExtra(eventId)
		}
		hdr := [1]uint64{1}
		// Note: write "knows" that the argument is &gp.labels,
		// because otherwise its write barrier behavior may not
		// be correct. See the long comment there before
		// changing the argument here.
		pmuprof[eventId].log.write(&gp.labels, nanotime(), hdr[:], stk)
	}

	atomic.Store(&profs[eventId].signalLock, 0)
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
func (p *pmuProfile) addNonGo(stk []uintptr, eventId int) {
	// (Other calls to add or addNonGo should be blocked out
	// by the fact that only one SIGPROF can be handled by the
	// process at a time. If not, this lock will serialize those too.)
	for !atomic.Cas(&profs[eventId].signalLock, 0, 1) {
		osyield()
	}

	if pmuprof[eventId].numExtra+1+len(stk) < len(pmuprof[eventId].extra) {
		i := pmuprof[eventId].numExtra
		pmuprof[eventId].extra[i] = uintptr(1 + len(stk))
		copy(pmuprof[eventId].extra[i+1:], stk)
		pmuprof[eventId].numExtra += 1 + len(stk)
	} else {
		pmuprof[eventId].lostExtra++
	}

	atomic.Store(&profs[eventId].signalLock, 0)
}

// addExtra adds the "extra" profiling events,
// queued by addNonGo, to the profile log.
// addExtra is called either from a signal handler on a Go thread
// or from an ordinary goroutine; either way it can use stack
// and has a g. The world may be stopped, though.
func (p *pmuProfile) addExtra(eventId int) {
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
		pmuprof[eventId].log.write(nil, 0, hdr[:], lostStk[:])
		p.lostExtra = 0
	}
}

func (p *pmuProfile) addLostAtomic64(count uint64, eventId int) {
	hdr := [1]uint64{count}
	lostStk := [2]uintptr{
		funcPC(_LostSIGPROFDuringAtomic64) + sys.PCQuantum,
		funcPC(_System) + sys.PCQuantum,
	}
	pmuprof[eventId].log.write(nil, 0, hdr[:], lostStk[:])
}

//go:linkname runtime_pprof_readPMUProfile runtime/pprof.readPMUProfile
func runtime_pprof_readPMUProfile(eventId int) ([]uint64, []unsafe.Pointer, bool) {
	lock(&pmuprof[eventId].lock)
	log := pmuprof[eventId].log
	unlock(&pmuprof[eventId].lock)
	data, tags, eof := log.read(profBufBlocking)
	if len(data) == 0 && eof {
		lock(&pmuprof[eventId].lock)
		pmuprof[eventId].log = nil
		unlock(&pmuprof[eventId].lock)
	}
	return data, tags, eof
}
