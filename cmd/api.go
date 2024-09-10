package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

type Entries struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Url      string  `json:"url"`
	Duration float32 `json:"duration"`
	Channel  string  `json:"channel"`
}

type playlist struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Entries     []Entries `json:"entries"`
}

// Download playlist JSON data from YouTube
func FetchPlaylist(id string) playlist {
	fmt.Println("Fetching playlist: " + id)
	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", fmt.Sprintf("https://www.youtube.com/playlist?list=%s", id))

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

func PercentageOf(total, percent int) int {
	return (total * percent) / 100
}

func LoadPlaylistCached(id string) (playlist, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	filePath := usr.HomeDir + "/.cache/ytm-tui/" + id + ".json"
	p := playlist{}

	file, err := os.Open(filePath)
	if err != nil {
		return playlist{}, fmt.Errorf("playlist not cached: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&p); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	return p, nil
}

func WriteToCache(id string, bytes []byte) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	filePath := usr.HomeDir + "/.cache/ytm-tui/" + id + ".json"

	err = os.MkdirAll(usr.HomeDir+"/.cache/ytm-tui", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(json.RawMessage(bytes)); err != nil {
		log.Fatalf("failed to encode JSON: %v", err)
	}

	log.Println("Data successfully written to cache.")
}

func QuickLoadPlaylist(id string) playlist {
	p, err := LoadPlaylistCached(id)
	if err != nil {
		p = FetchPlaylist(id)
	}
	return p
}
