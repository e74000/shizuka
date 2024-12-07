package cmd

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/e74000/shizuka/shizuka"
	"os"
)

const (
	ConfigPath = "shizuka_conf.json"

	DefaultSrc  = "site"
	DefaultDst  = "dist"
	DegaultPort = "8080"
)

// Config represents the structure of shizuka_conf.json
type Config struct {
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Port string `json:"port"`
}

// GetConfig loads the configuration from shizuka_conf.json or returns default values.
func GetConfig() Config {
	// Default configuration
	defaultConfig := Config{
		Src:  DefaultSrc,
		Dst:  DefaultDst,
		Port: DegaultPort,
	}

	// Open shizuka_conf.json
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Debug("failed to open config, using default values", "err", err)
		return defaultConfig
	}
	defer file.Close()

	// Decode JSON
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Error("failed to parse config, using default values", "err", err)
		return defaultConfig
	}

	// Fill missing values with defaults
	if config.Src == "" {
		config.Src = defaultConfig.Src
	}
	if config.Dst == "" {
		config.Dst = defaultConfig.Dst
	}
	if config.Port == "" {
		config.Port = defaultConfig.Port
	}

	return config
}

// WriteConfig writes the given configuration to shizuka_conf.json.
func WriteConfig(config Config) error {
	// Open the config file for writing, create it if it doesn't exist
	file, err := os.Create(ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON encoder
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print with indentation for readability

	// Encode the config into JSON and write it to the file
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func buildSite(src, dst string, opts *shizuka.BuildOpts) error {
	pb := shizuka.NewPageBuilder(src, dst)
	if opts != nil {
		pb.Opts = *opts
	}

	if err := pb.Index(); err != nil {
		return err
	}

	if err := pb.Build(); err != nil {
		return err
	}

	return nil
}
