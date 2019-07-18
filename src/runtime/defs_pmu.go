package runtime

const (
	GO_COUNT_HW_CPU_CYCLES              = 0x0
	GO_COUNT_HW_INSTRUCTIONS            = 0x1
	GO_COUNT_HW_CACHE_REFERENCES        = 0x2
	GO_COUNT_HW_CACHE_MISSES            = 0x3
	GO_COUNT_HW_RAW			    = 0x4
	GO_COUNT_PMU_EVENTS_MAX		    = 0x5
)

type PMUEventAttr struct {
	Period           uint64
	RawEvent         uint64
	PreciseIP        uint8
	IsKernelIncluded bool
	IsHvIncluded     bool
}
