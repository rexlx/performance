package performance

import "testing"

func TestMemory(t *testing.T) {
	m := GetMemoryUsage()
	t.Log(m.String())
}
