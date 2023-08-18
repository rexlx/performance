package performance

import "testing"

func TestMemory(t *testing.T) {
	m := GetMemoryUsage()
	t.Log(m.String())
}

func TestBytes(t *testing.T) {
	t.Log(Bytes(KiB * 6186))
	t.Log(ByteConverter(KiB * 6669))
}
