package runtime

// These Constants are platform agnostic and exposed to pprof.
// We use perfEventOpt to map these to the underlying OS and HW.
const (
	GO_COUNT_HW_CPU_CYCLES             = 0x0
	GO_COUNT_HW_INSTRUCTIONS           = 0x1
	GO_COUNT_HW_CACHE_REFERENCES       = 0x2
	GO_COUNT_HW_CACHE_MISSES           = 0x3
	GO_COUNT_HW_CACHE_LL_READ_ACCESSES = 0x4
	GO_COUNT_HW_CACHE_LL_READ_MISSES   = 0x5
	GO_COUNT_HW_RAW                    = 0x6
	GO_COUNT_PMU_EVENTS_MAX            = 0x7
)

type PMUEventAttr struct {
	Period                    uint64
	RawEvent                  uint64
	PreciseIP                 uint8
	IsSampleIPIncluded        bool
	IsSampleThreadIDIncluded  bool
	IsSampleAddrIncluded	  bool
	IsSampleCallchainIncluded bool
	IsKernelIncluded          bool
	IsHvIncluded              bool
	IsIdleIncluded            bool
	IsCallchainKernelIncluded bool
	IsCallchainUserIncluded   bool
}
