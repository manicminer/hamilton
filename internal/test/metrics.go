package test

import (
	"fmt"
	"runtime"
	"runtime/metrics"
	"time"
)

// WithMetrics runs f and prints the total CPU time and max heap memory used during its execution
func WithMetrics(f func()) {
	cpuSample := make([]metrics.Sample, 1)
	cpuSample[0].Name = "/cpu/classes/user:cpu-seconds"
	memorySample := make([]metrics.Sample, 1)
	memorySample[0].Name = "/memory/classes/heap/objects:bytes"
	runtime.GC()

	ticker := time.NewTicker(10 * time.Millisecond)
	done := make(chan bool)
	go func() {
		maxMemory := uint64(0)
		for {
			select {
			case <-done:
				fmt.Printf("Max heap memory: %d MiB\n", maxMemory/1024/1024)
				return
			case _ = <-ticker.C:
				metrics.Read(memorySample)
				memory := memorySample[0].Value.Uint64()
				if memory > maxMemory {
					maxMemory = memory
				}
			}
		}
	}()

	metrics.Read(cpuSample)
	cpuSecondsBefore := cpuSample[0].Value.Float64()

	f()

	metrics.Read(cpuSample)
	cpuSecondsAfter := cpuSample[0].Value.Float64()
	fmt.Printf("Total CPU time: %.0f ms\n", (cpuSecondsAfter-cpuSecondsBefore)*1000)

	ticker.Stop()
	done <- true
}
