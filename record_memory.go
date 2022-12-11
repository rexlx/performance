package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type HumanReadableBytes struct {
	Value  float64
	Symbol string
}

type MemoryStats struct {
	SwapTotal,
	SwapFree,
	MemTotal,
	MemFree,
	Buffers,
	Cached,
	Slab int
}

type MemoryUsage struct {
	Time time.Time
	Swap,
	Used,
	Total,
	Free,
	Buff,
	Cached,
	Slab float64
}

func (m *MemoryStats) analyzeUsage() *MemoryUsage {
	return &MemoryUsage{
		Time:   time.Now(),
		Swap:   float64(m.SwapTotal) - float64(m.SwapFree),
		Used:   float64(m.MemTotal) - float64(m.MemFree) - float64(m.Buffers) - float64(m.Cached) - float64(m.Slab),
		Free:   float64(m.MemFree),
		Total:  float64(m.MemTotal),
		Buff:   float64(m.Buffers),
		Cached: float64(m.Cached),
		Slab:   float64(m.Slab),
	}
}

func convertBytes(n float64) *HumanReadableBytes {
	h := &HumanReadableBytes{}
	for _, i := range []string{"KiB", "MiB", "GiB", "TiB"} {
		if n < 1024 || i == "TiB" {
			break
		}
		n /= 1024.0
		h.Value = n
		h.Symbol = i
	}
	return h
}

// GetMemoryUsage reads the system memory stats and calculates the usage
func GetMemoryUsage() *MemoryUsage {
	contents, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Println(err)
	}
	lines := strings.Split(string(contents), "\n")
	stats := generateMemoryMap(lines)
	memStats := &MemoryStats{
		Buffers:   stats["Buffers"],
		Cached:    stats["Cached"],
		MemTotal:  stats["MemTotal"],
		MemFree:   stats["MemFree"],
		SwapTotal: stats["SwapTotal"],
		SwapFree:  stats["SwapFree"],
		Slab:      stats["Slab"],
	}
	return memStats.analyzeUsage()
}

func generateMemoryMap(lines []string) map[string]int {
	memoryMap := make(map[string]int)
	for _, line := range lines {
		var i int
		data := strings.Split(line, ":")
		if len(data) < 2 {
			continue
		}
		if _, err := fmt.Sscanf(data[1], "%v kB", &i); err == nil {
			memoryMap[data[0]] = i
		}

	}
	return memoryMap
}
