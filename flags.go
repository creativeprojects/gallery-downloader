package main

import (
	"flag"
	"gallery-downloader/scan"
	"strings"
)

// Flags from command line
type Flags struct {
	ConfigFile  string
	Source      string
	Base        string
	Type        string
	Output      string
	Referer     string
	User        string
	Password    string
	WaitMin     int
	WaitMax     int
	InsecureTLS bool
}

func loadFlags() Flags {
	flags := Flags{}
	flag.StringVar(&flags.ConfigFile, "config", "config.json", "configuration file")
	flag.StringVar(&flags.Source, "source", "", "source HTML gallery")
	flag.StringVar(&flags.Base, "base", "", "base URL when downloading relative images")
	flag.StringVar(&flags.Type, "type", scan.AvailableGalleryScanners[0], "type of gallery ("+strings.Join(scan.AvailableGalleryScanners[:], ", ")+")")
	flag.StringVar(&flags.Output, "output", "", "output folder to store pictures")
	flag.StringVar(&flags.Referer, "referer", "", "referer header for HTML file, or for downloading images from a local HTML file")
	flag.StringVar(&flags.User, "user", "", "user (if the http server needs basic authentication)")
	flag.StringVar(&flags.Password, "password", "", "password (if the http server needs basic authentication)")
	flag.IntVar(&flags.WaitMin, "min-wait", 1000, "wait n milliseconds minimum before downloading the next image. Use 0 to deactivate")
	flag.IntVar(&flags.WaitMax, "max-wait", 3000, "wait n milliseconds maximum before downloading the next image. Use 0 to deactivate")
	flag.BoolVar(&flags.InsecureTLS, "insecure-tls", false, "Skip TLS certificate verification. Should only be enabled for testing locally")
	flag.Parse()
	return flags
}
