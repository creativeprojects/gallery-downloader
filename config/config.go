package config

import (
	"encoding/json"
	"io"
	"os"
)

// Configuration contains all configuration from JSON file
type Configuration struct {
	Browser Browser `json:"browser"`
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
