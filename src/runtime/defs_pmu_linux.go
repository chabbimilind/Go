package runtime

var perfEventOpt = []struct {
	Type    uint32 // type of event
	Config  uint64 // event 
} {
	{PERF_TYPE_HARDWARE, PERF_COUNT_HW_CPU_CYCLES},	      // index: GO_COUNT_HW_CPU_CYCLES
	{PERF_TYPE_HARDWARE, PERF_COUNT_HW_INSTRUCTIONS},     // index: GO_COUNT_HW_INSTRUCTIONS
	{PERF_TYPE_HARDWARE, PERF_COUNT_HW_CACHE_REFERENCES}, // index: GO_COUNT_HW_CACHE_REFERENCES
	{PERF_TYPE_HARDWARE, PERF_COUNT_HW_CACHE_MISSES},     // index: GO_COUNT_HW_CACHE_MISSES
	// TODO: add more perf events
}
