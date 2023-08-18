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
	Swap   int       `json:"swap_used"`
	Used   int       `json:"percent_used"`
	Total  int       `json:"total_memory"`
	Free   int       `json:"free"`
	Buff   int       `json:"buffered"`
	Cached int       `json:"cached"`
	Slab   int       `json:"slab"`
}

// analyzeUsage does some basic maths to determine the usage.
func (m *MemoryStats) analyzeUsage() *MemoryUsage {
	return &MemoryUsage{
		Time:   time.Now(),
		Swap:   m.SwapTotal - m.SwapFree,
		Used:   m.MemTotal - m.MemFree - m.Buffers - m.Cached - m.Slab,
		Free:   m.MemFree,
		Total:  m.MemTotal,
		Buff:   m.Buffers,
		Cached: m.Cached,
		Slab:   m.Slab,
	}
}

func (m *MemoryUsage) String() string {
	return fmt.Sprintf(
		memoryUsageTemplate,
		m.Time,
		Bytes(m.Used),
		Bytes(m.Swap),
		Bytes(m.Total),
		Bytes(m.Free),
		Bytes(m.Cached))
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

// define a template for the output. we need to pad for the floats
const memoryUsageTemplate = `Memory Usage
=====================
Time: %v
=====================
Used: %s Swap: %s Total: %s Free: %s Cached: %s
`

func Bytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
