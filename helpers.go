package performance

import "fmt"

const (
	_ = 1 << (iota * 10)
	KiB
	MiB
	GiB
	TiB
)

func ByteConverter(n int) string {
	switch {
	case n > TiB:
		return fmt.Sprintf("%.2f TiB", float64(n/TiB))
	case n > GiB:
		return fmt.Sprintf("%.2f GiB", float64(n/GiB))
	case n > MiB:
		return fmt.Sprintf("%.2f MiB", float64(n/MiB))
	case n > KiB:
		return fmt.Sprintf("%.2f KiB", float64(n/KiB))
	default:
		return fmt.Sprintf("%.2f   B", float64(n))
	}
}
