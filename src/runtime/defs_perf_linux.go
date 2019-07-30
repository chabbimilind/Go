package runtime

// Convert platform-agnostic pmu events to Linux perf events
var perfEventOpt = map[int32]struct {
	Type    uint32 // type of event
	Config  uint64 // event
} {
	GO_COUNT_HW_CPU_CYCLES : {_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CPU_CYCLES},
	GO_COUNT_HW_INSTRUCTIONS: {_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_INSTRUCTIONS},
	GO_COUNT_HW_CACHE_REFERENCES: {_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CACHE_REFERENCES},
	GO_COUNT_HW_CACHE_MISSES: {_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CACHE_MISSES},
	GO_COUNT_HW_CACHE_LL_READ_ACCESSES: {_PERF_TYPE_HW_CACHE, (_PERF_COUNT_HW_CACHE_LL) | (_PERF_COUNT_HW_CACHE_OP_READ << 8) | (_PERF_COUNT_HW_CACHE_RESULT_ACCESS << 16)},
	GO_COUNT_HW_CACHE_LL_READ_MISSES: {_PERF_TYPE_HW_CACHE, (_PERF_COUNT_HW_CACHE_LL) | (_PERF_COUNT_HW_CACHE_OP_READ << 8) | (_PERF_COUNT_HW_CACHE_RESULT_MISS << 16)},
	GO_COUNT_HW_RAW: {_PERF_TYPE_RAW, 0 /* will not be used */},
	// TODO: add more perf events
}

type perfEventAttr struct {
	Type               uint32
	Size               uint32
	Config             uint64
	Sample             uint64
	Sample_type        uint64
	Read_format        uint64
	Bits               uint64
	Wakeup             uint32
	Bp_type            uint32
	Ext1               uint64
	Ext2               uint64
	Branch_sample_type uint64
	Sample_regs_user   uint64
	Sample_stack_user  uint32
	Clockid            int32
	Sample_regs_intr   uint64
	Aux_watermark      uint32
	Sample_max_stack   uint16
	_                  uint16
}
