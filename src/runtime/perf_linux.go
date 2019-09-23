package runtime

import (
	"unsafe"
)

func ioctl(fd int32, req, arg int) (r, err int)
//go:noescape
func perfEventOpen(attr *perfEventAttr, pid, cpu, groupFd, flags, dummy int) (r int32, r2, err int)

const perfDataPages = 2
var perfPageSize uint64
var perfPageMask uint64

func perfAttrInit(eventId int32, eventAttr *PMUEventAttr, perfAttr *perfEventAttr) {
	perfAttr.Type = perfEventOpt[eventId].Type
	perfAttr.size = uint32(unsafe.Sizeof(*perfAttr))

	if eventId == GO_COUNT_HW_RAW {
		perfAttr.config = eventAttr.RawEvent
	} else {
		perfAttr.config = perfEventOpt[eventId].config
	}

	perfAttr.sample = eventAttr.Period
	if perfAttr.sample == 0 {
		perfAttr.read_format = _PERF_FORMAT_TOTAL_TIME_ENABLED | _PERF_FORMAT_TOTAL_TIME_RUNNING
	}

	if eventAttr.IsSampleIPIncluded {
		perfAttr.sample_type = _PERF_SAMPLE_IP
	}
	if eventAttr.IsSampleAddrIncluded {
		perfAttr.sample_type |= _PERF_SAMPLE_ADDR
	}
	if eventAttr.IsSampleCallchainIncluded {
		perfAttr.sample_type |= _PERF_SAMPLE_CALLCHAIN
	}
	if eventAttr.IsSampleThreadIDIncluded {
		perfAttr.sample_type |= _PERF_SAMPLE_TID
	}

	perfAttr.bits = 1 // the counter is disabled and will be enabled later
	perfAttr.bits |= uint64(eventAttr.PreciseIP) << 15 // precise ip
	if !eventAttr.IsKernelIncluded { // don't count kernel
		perfAttr.bits |= 1 << 5
	}
	if !eventAttr.IsHvIncluded { // don't count hypervisor
		perfAttr.bits |= 1 << 6
	}
	if !eventAttr.IsIdleIncluded { // don't count when idle
		perfAttr.bits |= 1 << 7
	}
	if !eventAttr.IsCallchainKernelIncluded {
		perfAttr.bits |= 1 << 21
	}
	if !eventAttr.IsCallchainUserIncluded {
		perfAttr.bits |= 1 << 22
	}
}

func perfMmapInit() {
	perfPageSize = uint64(physPageSize)
	perfPageMask = perfPageSize * perfDataPages - 1
}

func perfMmapSize() uintptr {
	if perfPageSize == 0 {
		println("The perf page size has been unknown!")
	}
	return uintptr(perfPageSize * (perfDataPages + /* metadata page */ 1))
}

func perfSetMmap(fd int32) *perfEventMmapPage {
	if perfPageSize == 0 {
		perfMmapInit()
	}

	size := perfMmapSize()
	r, err := mmap(nil, size, _PROT_WRITE | _PROT_READ, _MAP_SHARED, fd, /* page offset */ 0)
	if err != 0 {
		return nil
        }

	perfMmap := (*perfEventMmapPage)(r)
	// There is no memset available in Go runime
	// We instead employ the following approach to initialize all bytes in perfMmap to zeros
	var temp perfEventMmapPage
	memmove(unsafe.Pointer(perfMmap), unsafe.Pointer(&temp), unsafe.Sizeof(temp))

	return perfMmap
}

func perfUnsetMmap(mmapBuf *perfEventMmapPage) {
	size := perfMmapSize()
	munmap(unsafe.Pointer(mmapBuf), size)
}

func perfSkipNBytes(mmapBuf *perfEventMmapPage, n uint64 ) {
	head := mmapBuf.data_head
	tail := mmapBuf.data_tail
	if tail + n > head {
		n = head - tail
	}
	mmapBuf.data_tail += n
}

func perfSkipRecord(mmapBuf *perfEventMmapPage, hdr *perfEventHeader) {
	if mmapBuf == nil {
		return
	}
	remains := uint64(hdr.size) - uint64(unsafe.Sizeof(*hdr))
	if remains > 0 {
		perfSkipNBytes(mmapBuf, remains)
	}
}

func perfSkipAll(mmapBuf *perfEventMmapPage) {
	if mmapBuf == nil {
		return
	}
	remains := mmapBuf.data_head - mmapBuf.data_tail
	if remains > 0 {
		perfSkipNBytes(mmapBuf, remains)
	}
}

func perfReadNbytes(mmapBuf *perfEventMmapPage, buf unsafe.Pointer, n uint64) bool {
	if mmapBuf == nil {
		return false
	}

	// front of the circular data buffer	
	data := unsafe.Pointer(uintptr(unsafe.Pointer(mmapBuf)) + uintptr(perfPageSize))

	tail := mmapBuf.data_tail
	head := mmapBuf.data_head

	// compute bytes available in the circular buffer
	byteAvailable := head - tail
	if n > byteAvailable {
		return false
	}

	// compute offset of tail in the circular buffer
	tail &= perfPageMask

	bytesAtRight := (perfPageMask + 1) - tail

	// bytes to copy to the right of tail
	var right uint64
	if bytesAtRight < n {
		right = bytesAtRight
	} else {
		right = n
	}

	// copy bytes from tail position 
	memmove(buf, unsafe.Pointer(uintptr(data) + uintptr(tail)), uintptr(right))

	// if necessary, wrap and continue copy from left edge of buf
	if n > right {
		left := n - right
		memmove(unsafe.Pointer(uintptr(buf) + uintptr(right)) , data, uintptr(left))
	}

	// update tail after consuming n bytes
	mmapBuf.data_tail += n

	return true
}

func perfReadHeader(mmapBuf *perfEventMmapPage, hdr *perfEventHeader) bool {
	return perfReadNbytes(mmapBuf, unsafe.Pointer(hdr), uint64(unsafe.Sizeof(*hdr)))
}

func perfRecordSample(mmapBuf *perfEventMmapPage, eventAttr *PMUEventAttr, sampleData *perfSampleData) {
	if eventAttr.IsSampleIPIncluded {
		perfReadNbytes(mmapBuf, unsafe.Pointer(&(sampleData.ip)), uint64(unsafe.Sizeof(sampleData.ip)))
	}
	if eventAttr.IsSampleAddrIncluded {
		perfReadNbytes(mmapBuf, unsafe.Pointer(&(sampleData.addr)), uint64(unsafe.Sizeof(sampleData.addr)))
	}
	if eventAttr.IsSampleThreadIDIncluded {
		perfReadNbytes(mmapBuf, unsafe.Pointer(&(sampleData.pid)), uint64(unsafe.Sizeof(sampleData.pid)))
		perfReadNbytes(mmapBuf, unsafe.Pointer(&(sampleData.tid)), uint64(unsafe.Sizeof(sampleData.tid)))
	}
}

func perfStartCounter(fd int32) bool {
	_, err := ioctl(fd, _PERF_EVENT_IOC_ENABLE, 0)
	if err != 0 {
		println("Failed to enable the event count")
		return false
	}
	return true
}

func perfStopCounter(fd int32) bool {
	_, err := ioctl(fd, _PERF_EVENT_IOC_DISABLE, 0)
	if err != 0 {
		println("Failed to disable the event count")
		return false
	}
	return true
}

func perfResetCounter(fd int32) bool {
	_, err := ioctl(fd, _PERF_EVENT_IOC_RESET, 0)
	if err != 0 {
		println("Failed to reset the event count")
		return false
	}
	return true
}

func perfReadCounter(fd int32, val *uint64) bool {
	return read(fd, unsafe.Pointer(val), int32(unsafe.Sizeof(*val))) != -1
}
