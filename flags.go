package main

import (
	"flag"
	"strings"
)

// Configuration from command line
type Configuration struct {
	Source   string
	Base     string
	Type     string
	Output   string
	Agent    string
	User     string
	Password string
}

const (
	defaultAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:74.0) Gecko/20100101 Firefox/74.0"
)

var (
	configuration Configuration
	galleryTypes  = []string{
		"wowslider",
	}
)

func init() {
	flag.StringVar(&configuration.Source, "source", "", "source html")
	flag.StringVar(&configuration.Base, "base", "", "base URL when downloading relative images")
	flag.StringVar(&configuration.Type, "type", "", "type of gallery ("+strings.Join(galleryTypes, ", ")+")")
	flag.StringVar(&configuration.Output, "output", "", "output folder to store pictures")
	flag.StringVar(&configuration.Agent, "agent", defaultAgent, "browser user-agent")
	flag.StringVar(&configuration.User, "user", "", "user (for html loader only)")
	flag.StringVar(&configuration.Password, "password", "", "password (for html loader only)")
}
