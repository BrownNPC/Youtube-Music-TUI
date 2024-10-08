package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
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
