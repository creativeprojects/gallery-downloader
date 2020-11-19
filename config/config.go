package config

import (
	"encoding/json"
	"io"
	"os"
)

// Configuration contains all configuration from JSON file
type Configuration struct {
	Browser  Browser   `json:"browser"`
	Profiles []Profile `json:"profiles"`
}

// Browser contains all browser configuration
type Browser struct {
	Default Group `json:"default"`
	HTTP    Group `json:"http"`
	HTTPS   Group `json:"https"`
	HTTP2   Group `json:"http2"`
	HTML    Group `json:"html"`
	Picture Group `json:"picture"`
}

// Group contains browser configuration for each element (html, picture, etc.)
type Group struct {
	Headers map[string]string `json:"headers"`
}

// Profile contains the type of gallery and how to parse the images
type Profile struct {
	Priority      int    `json:"priority"`
	Name          string `json:"name"`
	Generator     string `json:"generator"`
	DetectGallery string `json:"detectGallery"`
	DetectImage   string `json:"detectImage"`
	MinWait       int    `json:"minWait"`
	MaxWait       int    `json:"maxWait"`
	Parallel      int    `json:"parallel"`
}

// newConfiguration creates an empty configuration object
func newConfiguration() *Configuration {
	return &Configuration{}
}

// LoadFileConfiguration loads a Configuration object from a JSON file
func LoadFileConfiguration(fileName string) (*Configuration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return loadConfiguration(file)
}

func loadConfiguration(reader io.ReadCloser) (*Configuration, error) {
	defer reader.Close()
	decoder := json.NewDecoder(reader)
	cfg := newConfiguration()
	err := decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
