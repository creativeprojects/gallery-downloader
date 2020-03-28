package main

import (
	"encoding/json"
	"io"
	"os"
)

// Configuration contains all configuration from JSON file
type Configuration struct {
	Browser BrowserConfiguration `json:"browser"`
}

// BrowserConfiguration contains all browser configuration
type BrowserConfiguration struct {
	UserAgent string               `json:"userAgent"`
	HTML      ElementConfiguration `json:"html"`
	Picture   ElementConfiguration `json:"picture"`
}

// ElementConfiguration contains browser configuration for each element (html, picture, etc.)
type ElementConfiguration struct {
	Headers map[string]string `json:"headers"`
}

// newConfiguration creates an empty configuration object
func newConfiguration() *Configuration {
	return &Configuration{}
}

func loadFileConfiguration(fileName string) (*Configuration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return loadConfiguration(file)
}

func loadConfiguration(reader io.ReadCloser) (*Configuration, error) {
	defer reader.Close()
	decoder := json.NewDecoder(reader)
	config := newConfiguration()
	err := decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
