package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
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

// Download playlist json data from youtube
func FetchPlaylist(id string) playlist {
	fmt.Println("Fetching playlist: " + id)
	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", fmt.Sprint("https://www.youtube.com/playlist?list=", id))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stderr: %s\n", stderr.String())
		os.Exit(1)
	}
	p := playlist{}
	json.Unmarshal(stdout.Bytes(), &p)
	// write to cache file in home/.cache/ytm-tui/
	WriteToCache(id, stdout.Bytes())
	return p
}
func PercentageOf(total, percent int) int {
	return (total * percent) / 100
}

func LoadPlaylistCached(id string) (playlist, error) {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	// Define the file path
	filePath := usr.HomeDir + "/.cache/ytm-tui/" + id + ".json"
	p := playlist{}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return playlist{}, fmt.Errorf("playlist not cached%v", err)
	}
	defer file.Close()

	// Decode the playlist from the file
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&p); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	return p, nil
}

func WriteToCache(id string, bytes []byte) {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	// Define the file path
	filePath := usr.HomeDir + "/.cache/ytm-tui/" + id + ".json"

	// Ensure the directory exists
	err = os.MkdirAll(usr.HomeDir+"/.cache/ytm-tui", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	// Encode the playlist to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: format the JSON with indentation
	if err := encoder.Encode(bytes); err != nil {
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
