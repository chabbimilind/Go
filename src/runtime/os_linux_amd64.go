// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/atomic"
	"unsafe"
)

//go:noescape
func perfEventOpen(attr *perfEventAttr, pid, cpu, groupFd, flags, dummy int64) (r int32, r2, err int64)

func ioctl(fd int32, req, arg int64) int64

// func fcntl(fd, cmd int64, arg interface{}) (r int64, err int64)
func fcntl(fd int32, cmd, arg int64) (r int64, err int64)

//go:noescape
func fcntl2(fd int32, cmd int64, arg *fOwnerEx) (r int64, err int64)

func setProcessPMUProfiler(eventAttr *PMUEventAttr) {
	if eventAttr != nil {
		// Enable the Go signal handler if not enabled.
		if atomic.Cas(&handlingSig[_SIGPMU], 0, 1) {
			atomic.Storeuintptr(&fwdSig[_SIGPMU], getsig(_SIGPMU))
			setsig(_SIGPMU, funcPC(sighandler))
		}
	} else {
		// If the Go signal handler should be disabled by default,
		// disable it if it is enabled.
		if !sigInstallGoHandler(_SIGPMU) {
			if atomic.Cas(&handlingSig[_SIGPMU], 1, 0) {
				setsig(_SIGPMU, atomic.Loaduintptr(&fwdSig[_SIGPMU]))
			}
		}
	}
}

func setThreadPMUProfiler(eventId int32, eventAttr *PMUEventAttr) {
	_g_ := getg()

	if eventAttr == nil {
		if _g_.m.eventAttrs[eventId] != nil {
			closefd(_g_.m.eventFds[eventId])
		}
	} else {
		var perfAttr perfEventAttr
		perfAttr.Size = uint32(unsafe.Sizeof(perfAttr))
		perfAttr.Type = perfEventOpt[eventId].Type
		perfAttr.Config = perfEventOpt[eventId].Config
		perfAttr.Sample = eventAttr.Period
		perfAttr.Bits = uint64(eventAttr.PreciseIP) << 15 // precise ip
		if !eventAttr.IsKernelIncluded { // don't count kernel
			perfAttr.Bits += 0b100000
		}
		if !eventAttr.IsHvIncluded { // don't count hypervisor
			perfAttr.Bits += 0b1000000
		}

		fd, _, _ := perfEventOpen(&perfAttr, 0, -1, -1, 0, /* dummy */ 0)
		_g_.m.eventFds[eventId] = fd
		r, _ := fcntl(fd, /* F_GETFL */ 0x3, 0)
		fcntl(fd, /* F_SETFL */ 0x4, r | /* O_ASYNC */ 0x2000)
		fcntl(fd, /* F_SETSIG */ 0xa, _SIGPMU)
		fOwnEx := fOwnerEx{/* F_OWNER_TID */ 0, int32(gettid())}
		fcntl2(fd, /* F_SETOWN_EX */ 0xf, &fOwnEx)
	}

	_g_.m.eventAttrs[eventId] = eventAttr
}

//go:nowritebarrierrec
func sigpmuhandler(info *siginfo, ctxt unsafe.Pointer, gp *g) {
	fd := info.si_fd
	ioctl(fd, PERF_EVENT_IOC_DISABLE, 0)
	_g_ := getg()
	c := &sigctxt{info, ctxt}

	var eventId int = -1
	for i := 0; i < maxPMUEvent; i++ {
		if _g_.m.eventFds[i] == fd {
			eventId = i
			break
		}
	}
	if eventId != -1 {
		sigpmu(c.sigpc(), c.sigsp(), c.siglr(), gp, _g_.m, eventId)
	}

	ioctl(fd, PERF_EVENT_IOC_ENABLE, 0)
	return
}
