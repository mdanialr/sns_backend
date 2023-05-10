package helper

import (
	"fmt"
	"math"
)

// BytesToHumanize convert given bytes in int64 to human-readable units such as
// 5MB, 244.21KB.
func BytesToHumanize(b int64) string {
	const unitSizeI = 1000
	const unitSizeF = float64(unitSizeI)

	if b < unitSizeI {
		return fmt.Sprintf("%dB", b)
	}
	if b < int64(math.Pow(unitSizeF, 2)) {
		return fmt.Sprintf("%.2fKB", float64(b)/unitSizeF)
	}
	if b < int64(math.Pow(unitSizeF, 3)) {
		return fmt.Sprintf("%.2fMB", float64(b)/unitSizeF/unitSizeF)
	}
	return fmt.Sprintf("%.2fGB", float64(b)/unitSizeF/unitSizeF/unitSizeF)
}
