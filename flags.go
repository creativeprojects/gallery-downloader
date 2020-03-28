package main

import (
	"flag"
	"strings"
)

// Flags from command line
type Flags struct {
	ConfigFile string
	Source     string
	Base       string
	Type       string
	Output     string
	Referer    string
	User       string
	Password   string
	Wait       int
}

var (
	galleryTypes = []string{
		"AnchorHREF",
	}
)

func loadFlags() Flags {
	flags := Flags{}
	flag.StringVar(&flags.ConfigFile, "config", "config.json", "configuration file")
	flag.StringVar(&flags.Source, "source", "", "source HTML gallery")
	flag.StringVar(&flags.Base, "base", "", "base URL when downloading relative images")
	flag.StringVar(&flags.Type, "type", galleryTypes[0], "type of gallery ("+strings.Join(galleryTypes, ", ")+")")
	flag.StringVar(&flags.Output, "output", "", "output folder to store pictures")
	flag.StringVar(&flags.Referer, "referer", "", "referer header for HTML file, or for downloading images from a local HTML file")
	flag.StringVar(&flags.User, "user", "", "user (if the http server needs basic authentication)")
	flag.StringVar(&flags.Password, "password", "", "password (if the http server needs basic authentication)")
	flag.IntVar(&flags.Wait, "wait", 3000, "wait between 1000 and n milliseconds before downloading the next image. Use 0 to deactivate")
	flag.Parse()
	return flags
}
