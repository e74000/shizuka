package cmd

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/e74000/shizuka/shizuka"
	"os"
)

const (
	ConfigPath = "shizuka_conf.json"
)

var DefaultConf = Config{
	Src:  "site",
	Dst:  "dist",
	Port: "8080",

	UseSitemap: false,
	UseRSS:     false,
	BaseURL:    "",
}

// Config represents the structure of shizuka_conf.json
type Config struct {
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Port string `json:"port"`

	UseSitemap bool `json:"use_sitemap"`
	UseRSS     bool `json:"use_rss"`

	BaseURL string `json:"base_url"`

	SiteTitle       string `json:"site_title"`
	SiteDescription string `json:"site_description"`
	SiteLang        string `json:"site_lang"`
}

// GetConfig loads the configuration from shizuka_conf.json or returns default values.
func GetConfig() Config {
	// Default configuration
	defaultConfig := Config{
		Src:        DefaultConf.Src,
		Dst:        DefaultConf.Dst,
		Port:       DefaultConf.Port,
		UseSitemap: DefaultConf.UseSitemap,
		BaseURL:    DefaultConf.BaseURL,
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
	if config.BaseURL == "" {
		config.BaseURL = defaultConfig.BaseURL
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

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func makeOpts(config Config) *shizuka.BuildOpts {
	return &shizuka.BuildOpts{
		UseSitemap:      config.UseSitemap,
		UseRss:          config.UseRSS,
		BaseURL:         config.BaseURL,
		SiteTitle:       config.SiteTitle,
		SiteDescription: config.SiteDescription,
		SiteLang:        config.SiteLang,
	}
}
