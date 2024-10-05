package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// getCacheDir returns the .cache directory in the user's home directory
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	var cacheDir string
	switch runtime.GOOS {
	case "windows":
		// On Windows, we can use %APPDATA% or %LOCALAPPDATA% for cache directories, but we’ll stick with the home directory.
		cacheDir = filepath.Join(home, "AppData", "Local", ".cache", "ytt")
	default:
		// Linux/macOS
		cacheDir = filepath.Join(home, ".cache", "ytt")
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create .cache directory: %w", err)
	}

	return cacheDir, nil
}

// downloadYTDLP downloads yt-dlp to the .cache folder
func downloadYTDLP() (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// Define the download URL based on the operating system
	var url string
	switch runtime.GOOS {
	case "windows":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "darwin", "linux":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Set the destination path for the yt-dlp binary
	fileName := "yt-dlp"
	if runtime.GOOS == "windows" {
		fileName += ".exe"
	}
	filePath := filepath.Join(cacheDir, fileName)

	// Download yt-dlp binary
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download yt-dlp: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download yt-dlp: received status code %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create yt-dlp binary: %w", err)
	}
	defer out.Close()

	// Copy the content from the response to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save yt-dlp binary: %w", err)
	}

	// Make the binary executable on Unix-like systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filePath, 0755); err != nil {
			return "", fmt.Errorf("failed to make yt-dlp executable: %w", err)
		}
	}

	return filePath, nil
}
