package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

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

//go:embed defaultconfig.toml
var defaultConfig string

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
	_, err = file.WriteString(defaultConfig)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}
}
