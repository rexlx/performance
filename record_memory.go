package performance

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	SwapTotal = iota
	SwapFree
	MemTotal
	MemFree
	Buffers
	Cached
	Slab
)

// MemoryStats type holds raw data used to calculate usage
type MemoryStats struct {
	SwapTotal,
	SwapFree,
	MemTotal,
	MemFree,
	Buffers,
	Cached,
	Slab int
}

// MemoryUsage type represents the memory util at a given time
type MemoryUsage struct {
	Time   time.Time `json:"time"`
	Swap   float64   `json:"swap_used"`
	Used   float64   `json:"percent_used"`
	Total  float64   `json:"total_memory"`
	Free   float64   `json:"free"`
	Buff   float64   `json:"buffered"`
	Cached float64   `json:"cached"`
	Slab   float64   `json:"slab"`
}

// analyzeUsage does some basic maths to determine the usage.
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

// GetMemoryUsage reads the system memory stats and calculates the usage
func GetMemoryUsage() *MemoryUsage {
	contents, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Println(err)
	}
	lines := strings.Split(string(contents), "\n")
	stats := generateMemoryMap(lines)
	memStats := &MemoryStats{
		Buffers:   GetKeyOrNone("Buffers", stats),
		Cached:    GetKeyOrNone("Cached", stats),
		MemTotal:  GetKeyOrNone("MemTotal", stats),
		MemFree:   GetKeyOrNone("MemFree", stats),
		SwapTotal: GetKeyOrNone("SwapTotal", stats),
		SwapFree:  GetKeyOrNone("SwapFree", stats),
		Slab:      GetKeyOrNone("Slab", stats),
	}
	return memStats.analyzeUsage()
}

func GetKeyOrNone(key string, m map[string]int) int {
	if val, ok := m[key]; ok {
		return val
	}
	return 0
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
