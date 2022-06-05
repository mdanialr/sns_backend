package service

import "fmt"

// FormatBytesToHumanString format given bytes to human-readable size like B, KB, MB.
func FormatBytesToHumanString(size int64) string {
	const unit = 1024

	if size < unit {
		// Bytes | 1024
		return fmt.Sprintf("%dB", size)
	}
	if size >= unit && size < (unit*unit) {
		// KiloBytes | Should less than unit**2
		res := float64(size) / unit
		return fmt.Sprintf("%.1fKB", res)
	}
	if size >= (unit*unit) && size < (unit*unit*unit) {
		// MegaBytes | Should less than unit**3
		res := float64(size) / (unit * unit)
		return fmt.Sprintf("%.1fMB", res)
	}

	// GigaBytes | Should more than unit**3
	res := float64(size) / (unit * unit * unit)
	return fmt.Sprintf("%.1fGB", res)
}
