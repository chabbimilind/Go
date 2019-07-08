// created by cgo -cdefs and then converted to Go
// cgo -cdefs defs_linux.go defs1_linux.go

package runtime

const (
	GO_COUNT_HW_CPU_CYCLES              = 0x0
	GO_COUNT_HW_INSTRUCTIONS            = 0x1
	GO_COUNT_HW_CACHE_REFERENCES        = 0x2
	GO_COUNT_HW_CACHE_MISSES            = 0x3
	GO_COUNT_HW_BRANCH_INSTRUCTIONS     = 0x4
	GO_COUNT_HW_BRANCH_MISSES           = 0x5
	GO_COUNT_HW_BUS_CYCLES              = 0x6
	GO_COUNT_HW_STALLED_CYCLES_FRONTEND = 0x7
	GO_COUNT_HW_STALLED_CYCLES_BACKEND  = 0x8
	GO_COUNT_HW_REF_CPU_CYCLES          = 0x9

	GO_COUNT_HW_CACHE_L1D  = 0xa
	GO_COUNT_HW_CACHE_L1I  = 0xb
	GO_COUNT_HW_CACHE_LL   = 0xc
	GO_COUNT_HW_CACHE_DTLB = 0xd
	GO_COUNT_HW_CACHE_ITLB = 0xe
	GO_COUNT_HW_CACHE_BPU  = 0xf
	GO_COUNT_HW_CACHE_NODE = 0x10

	GO_COUNT_HW_CACHE_OP_READ     = 0x11
	GO_COUNT_HW_CACHE_OP_WRITE    = 0x12
	GO_COUNT_HW_CACHE_OP_PREFETCH = 0x13

	GO_COUNT_HW_CACHE_RESULT_ACCESS = 0x14
	GO_COUNT_HW_CACHE_RESULT_MISS   = 0x15

	GO_COUNT_SW_CPU_CLOCK        = 0x16
	GO_COUNT_SW_TASK_CLOCK       = 0x17
	GO_COUNT_SW_PAGE_FAULTS      = 0x18
	GO_COUNT_SW_CONTEXT_SWITCHES = 0x19
)

type PMUEventAttr struct {
    Period           uint64
    PreciseIP        uint8
    IsKernelIncluded bool
    IsHvIncluded     bool
}
