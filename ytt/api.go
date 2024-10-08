package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
)

// ===============================================================
// LOCAL YT-DLP API
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

//go:embed help.txt
var HelpMessage string

func handleCommandLineArgs() {

	if len(os.Args) == 1 {
		return
	}

	switch os.Args[1] {
	case "help", "--help", "-h":

		fmt.Println(HelpMessage)
		os.Exit(0)
	case "refresh", "--refresh", "-r":
		ClearCache()
	case "config", "--config", "-c":
		// open the config file folder, os independant
		openConfigFolder()

		os.Exit(0)
	default:
		fmt.Println("unknown command, use 'ytt help'")
		os.Exit(0)
	}
}
