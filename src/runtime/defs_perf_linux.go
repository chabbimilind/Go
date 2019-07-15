package runtime

var perfEventOpt = []struct {
	Type    uint32 // type of event
	Config  uint64 // event
} {
	{_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CPU_CYCLES},	// index: GO_COUNT_HW_CPU_CYCLES
	{_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_INSTRUCTIONS},     // index: GO_COUNT_HW_INSTRUCTIONS
	{_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CACHE_REFERENCES}, // index: GO_COUNT_HW_CACHE_REFERENCES
	{_PERF_TYPE_HARDWARE, _PERF_COUNT_HW_CACHE_MISSES},     // index: GO_COUNT_HW_CACHE_MISSES
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
