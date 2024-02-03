package go_http_utils

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	// %s is POD_NAME or hostname
	runtimeKeyUptimeSeconds         = "go.%s.uptime_seconds"
	runtimeKeyGoroutines            = "go.%s.goroutines"
	runtimeKeyThreads               = "go.%s.threads"
	runtimeKeyGCPauseMicroseconds   = "go.%s.gc_pause_microseconds"
	runtimeKeyMemoryAllocBytes      = "go.%s.mem_alloc_bytes"
	runtimeKeyMemoryAllocBytesTotal = "go.%s.mem_alloc_bytes_total"
	runtimeKeyMemorySysBytes        = "go.%s.mem_sys_bytes"
	runtimeKeyMemoryHeapAllocBytes  = "go.%s.mem_heap_alloc_bytes"
	runtimeMaxGCPausesTimings       = 10
)

var runtimeStarted = time.Now()

// ServeRuntimeMetrics sends golang runtime metrics.
func ServeRuntimeMetrics(ctx context.Context, m Metrics) {
	name := os.Getenv(envPodName)
	if len(name) == 0 {
		name, _ = os.Hostname()
	}
	if name == "" {
		return
	}

	collect := runtimeCollector(name, m)

	interval := envDuration(envRuntimeMetricsInterval, defaultRuntimeMetricsInterval)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}

		collect()
	}
}

func runtimeCollector(name string, metrics Metrics) func() {
	name = strings.Replace(name, ".", "_", -1)

	uptimeSeconds := fmt.Sprintf(runtimeKeyUptimeSeconds, name)
	goroutines := fmt.Sprintf(runtimeKeyGoroutines, name)
	threads := fmt.Sprintf(runtimeKeyThreads, name)
	gcPauseMicroseconds := fmt.Sprintf(runtimeKeyGCPauseMicroseconds, name)
	memoryAllocBytes := fmt.Sprintf(runtimeKeyMemoryAllocBytes, name)
	memoryAllocBytesTotal := fmt.Sprintf(runtimeKeyMemoryAllocBytesTotal, name)
	memorySysBytes := fmt.Sprintf(runtimeKeyMemorySysBytes, name)
	memoryHeapAllocBytes := fmt.Sprintf(runtimeKeyMemoryHeapAllocBytes, name)

	lastNumGC := int64(0)

	var statGC debug.GCStats
	var statMem runtime.MemStats

	return func() {
		// uptime
		metrics.Gauge(uptimeSeconds, time.Since(runtimeStarted).Seconds())

		// goroutines and threads
		metrics.Gauge(goroutines, runtime.NumGoroutine())
		n, _ := runtime.ThreadCreateProfile(nil)
		metrics.Gauge(threads, n)

		// GC stats
		debug.ReadGCStats(&statGC)
		cnt := int(statGC.NumGC - lastNumGC)
		if cnt > len(statGC.Pause) {
			cnt = len(statGC.Pause) // 256*2+3
		}
		for i := 0; i < cnt && i < runtimeMaxGCPausesTimings; i++ {
			metrics.Duration(gcPauseMicroseconds, statGC.Pause[i]*time.Microsecond) // micro
		}
		lastNumGC = statGC.NumGC

		// memory stats
		runtime.ReadMemStats(&statMem)
		metrics.Gauge(memoryAllocBytes, statMem.Alloc)
		metrics.Gauge(memoryAllocBytesTotal, statMem.TotalAlloc)
		metrics.Gauge(memorySysBytes, statMem.Sys)
		metrics.Gauge(memoryHeapAllocBytes, statMem.HeapAlloc)
	}
}
