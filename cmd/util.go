package main

import "fmt"

func PercentageOf(total, percent int) int {
	return (total * percent) / 100
}
func formatTime(ms int) string {
	seconds := ms / 1000
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
