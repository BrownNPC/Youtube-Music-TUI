package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

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

// =====================================================================
// UTILITIES
// =====================================================================
func openFolder(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
func openConfigFolder() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	// Determine platform-specific config folder path

	configFolder := filepath.Join(usr.HomeDir, ".config", "ytt")

	err = openFolder(configFolder)
	if err != nil {
		log.Fatalf("failed to open config folder: %v", err)
	}
}

// getCacheDir returns the .cache directory in the user's home directory
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	cacheDir := filepath.Join(home, ".cache", "ytt")
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create .cache directory: %w", err)
	}
	return cacheDir, nil
}
