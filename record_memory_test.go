package performance

import "testing"

func TestMemory(t *testing.T) {
	m := GetMemoryUsage()
	t.Log(m.String())
}

func TestBytes(t *testing.T) {
	t.Log(Bytes(6177786666))
	t.Log(ByteConverter(KiB * 6177786))
}
