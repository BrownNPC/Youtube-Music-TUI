package main

import "fmt"

func PercentageOf(total, percent int) int {
	return (total * percent) / 100
}

type IntType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Generic function to format time
func formatTime[T IntType](seconds T) string {
	// Convert the input to int for further calculations
	sec := int(seconds)
	hours := sec / 3600
	minutes := (sec % 3600) / 60
	sec = sec % 60
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, sec)
	}
	return fmt.Sprintf("%02d:%02d", minutes, sec)
}
