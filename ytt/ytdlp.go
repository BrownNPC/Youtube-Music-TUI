package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func getYoutubeStreamURL(youtubeURL string) (string, error) {
	// Using yt-dlp to extract the audio stream URL

	ytdlpPath, err := downloadYTDLP()
	if err != nil {
		log.Fatalf("failed to install yt-dlp: %v", err)
	}

	cmd := exec.Command(ytdlpPath, "--get-url", "-f", "bestaudio", youtubeURL)

	// Create pipes for capturing stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error getting StdoutPipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error getting StderrPipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting command: %v", err)
	}

	// Use a WaitGroup to wait for both stdout and stderr to be processed
	var wg sync.WaitGroup
	wg.Add(2)

	var outputBuffer bytes.Buffer

	// Process stdout in a goroutine
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)

		// Read the output in chunks to handle long lines
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatalf("Error reading stdout: %v", err)
				}
				break
			}
			outputBuffer.WriteString(line)

		}
	}()

	// Process stderr in a goroutine
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)

		// Read stderr in chunks to handle long lines
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatalf("Error reading stderr: %v", err)
				}
				break
			}
			log.Printf("Command stderr: %s", line) // Optional: print the stderr in real time
		}
	}()

	// Wait for both stdout and stderr goroutines to finish
	wg.Wait()

	// Wait for the command to finish execution
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	output := outputBuffer.String()
	return strings.TrimSpace(string(output)), nil
}

func FetchPlaylist(id string) playlist {
	ytdlpPath, err := downloadYTDLP()
	if err != nil {
		log.Fatalf("failed to install yt-dlp: %v", err)
	}

	fmt.Println("Fetching playlist: " + id)
	cmd := exec.Command(ytdlpPath, "--flat-playlist", "-J", fmt.Sprintf("https://www.youtube.com/playlist?list=%s", id))

	// Create pipes for capturing stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error getting StdoutPipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error getting StderrPipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting command: %v", err)
	}

	// Use a WaitGroup to wait for both stdout and stderr to be processed
	var wg sync.WaitGroup
	wg.Add(2)

	var outputBuffer bytes.Buffer

	// Process stdout in a goroutine
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)

		// Read the output in chunks to handle long lines
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatalf("Error reading stdout: %v", err)
				}
				break
			}
			outputBuffer.WriteString(line)

		}
	}()

	// Process stderr in a goroutine
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)

		// Read stderr in chunks to handle long lines
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatalf("Error reading stderr: %v", err)
				}
				break
			}
			log.Printf("Command stderr: %s", line) // Optional: print the stderr in real time
		}
	}()

	// Wait for both stdout and stderr goroutines to finish
	wg.Wait()

	// Wait for the command to finish execution
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	// Parse the collected JSON output
	p := playlist{}
	if err := json.Unmarshal(outputBuffer.Bytes(), &p); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Cache the result after fetching
	WriteToCache(id, outputBuffer.Bytes())

	fmt.Println("Command finished")
	return p
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

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	} else if errors.Is(err, os.ErrNotExist) {
	}

	// File does not exist, proceed to download yt-dlp binary
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
