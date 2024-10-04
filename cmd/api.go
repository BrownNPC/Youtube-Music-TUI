package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/skratchdot/open-golang/open"
)

// ===============================================================
// YT-DLP API
// ===============================================================

var playlistCacheDIR = "/.cache/ytt/lists/"

type Entry struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Url      string  `json:"url"`
	Duration float32 `json:"duration"`
	Channel  string  `json:"channel"`
}

type playlist struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Entries     []Entry `json:"entries"`
}

// Download playlist JSON data from YouTube

func FetchPlaylist(id string) playlist {
	ytdlpPath, err := Install(context.Background(), &InstallOptions{})
	if err != nil {
		log.Fatalf("failed to install yt-dlp: %v", err)
	}

	fmt.Println("Fetching playlist: " + id)
	cmd := exec.Command(ytdlpPath.Executable, "--flat-playlist", "-J", fmt.Sprintf("https://www.youtube.com/playlist?list=%s", id))

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

func LoadPlaylistCached(id string) (playlist, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	filePath := usr.HomeDir + playlistCacheDIR + id + ".json"
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

	filePath := usr.HomeDir + playlistCacheDIR + id + ".json"

	err = os.MkdirAll(usr.HomeDir+playlistCacheDIR, os.ModePerm)
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

func getYoutubeStreamURL(youtubeURL string) (string, error) {
	// Using yt-dlp to extract the audio stream URL
	ytdlpPath, err := Install(context.Background(), &InstallOptions{})
	if err != nil {
		log.Fatalf("failed to install yt-dlp: %v", err)
	}

	cmd := exec.Command(ytdlpPath.Executable, "-f", "bestaudio[ext=m4a]", "-g", youtubeURL)

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

// ============================================================
// CONFIG FILE API
// ============================================================

// load config file (toml)

type Config struct {
	IDs []string `toml:"playlists"`
}

func LoadConfig() (Config, error) {

	//create default config file, if it doesn't exist
	MakeDefaultConfig()

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	filePath := usr.HomeDir + "/.config/ytt/config.toml"

	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("config file not found: %v", err)
	}

	var cfg Config
	toml.NewDecoder(file).Decode(&cfg)

	defer file.Close()

	return cfg, nil

}

func ClearCache() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	err = os.RemoveAll(usr.HomeDir + playlistCacheDIR)
	if err != nil {
		log.Fatalf("failed to clear cache: %v", err)
	}
}

func MakeDefaultConfig() {
	// make the file if it does not exist
	// return if it exists
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %v", err)
	}

	filePath := usr.HomeDir + "/.config/ytt/config.toml"

	_, err = os.Stat(filePath)
	if err == nil {
		return // file exists
	}

	err = os.MkdirAll(usr.HomeDir+"/.config/ytt", os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	// write the default config

	defaultConfig := `
#PASTE YOUTUBE PLAYLIST ID'S HERE, PLEASE DONT FORGET THE COMMAS
playlists = [
	#synthwave radio
	"PLkcA3mJSVisBLbLhQ6ZnTCi9nGHTVUDaI",
	# minecraft ost
	"PLefKpFQ8Pvy5aCLAGHD8Zmzsdljos-t2l"
]

`

	_, err = file.WriteString(defaultConfig)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}
}

// =====================================================================
// UTILITIES
// =====================================================================

func handleCommandLineArgs() {

	if len(os.Args) == 1 {
		return
	}

	switch os.Args[1] {
	case "help", "--help", "-h":
		help_message := `
usage:	ytt [options]
options:
  -h, help, --help       Show this help message
  -c, config, --config   Open config file folder
  -r, refresh, --refresh Refresh the playlist cache
  `

		fmt.Println(help_message)
		os.Exit(0)
	case "refresh", "--refresh", "-r":
		ClearCache()
	case "config", "--config", "-c":
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("failed to get current user: %v", err)
		}
		open.Run(usr.HomeDir + "/.config/ytt/")
		os.Exit(0)
	default:
		fmt.Println("unknown command, use 'ytt help'")
		os.Exit(0)
	}
}

// =====================================================================
