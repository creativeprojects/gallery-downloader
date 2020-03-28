package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if configuration.Source == "" {
		flag.Usage()
		log.Fatal("\nError: missing HTML source (-source)")
	}

	if configuration.Type == "" {
		flag.Usage()
		log.Fatal("\nError: mising gallery type (-type)")
	}

	if configuration.Output == "" {
		flag.Usage()
		log.Fatal("\nError: missing output folder (-output)")
	}
	if stat, err := os.Stat(configuration.Output); err == nil || os.IsExist(err) {
		if !stat.IsDir() {
			log.Fatalf("Output '%s' exists but is not a directory", configuration.Output)
		}
	}
	if _, err := os.Stat(configuration.Output); os.IsNotExist(err) {
		err = os.MkdirAll(configuration.Output, 0755)
		if err != nil {
			log.Fatalf("Cannot create output directory: %v", err)
		}
	}

	var baseURL = &url.URL{}
	var err error
	if configuration.Base != "" {
		baseURL, err = url.Parse(configuration.Base)
		if err != nil {
			log.Fatal("Error: -base value is not a parsable URL")
		}
		if !baseURL.IsAbs() {
			log.Fatal("Error: -base value is not an absolute URL")
		}
	}

	sourceURL, err := url.Parse(configuration.Source)
	if err != nil {
		log.Fatalf("Error parsing source URL: %v", err)
	}
	if sourceURL.Scheme == "" {
		// Let's consider this is a file on disk
		sourcefile, err := os.Open(configuration.Source)
		if err != nil {
			log.Fatalf("Error: cannot open HTML source file: %v", err)
		}
		pictures, err := loadWowSliderSource(sourcefile)
		if err != nil {
			log.Fatalf("Error: cannot parse HTML source file: %v", err)
		}
		if pictures == nil || len(pictures) == 0 {
			log.Println("No picture found in the HTML source!")
		}
		total := len(pictures)
		log.Printf("Found %d picture(s) in the HTML source", total)

		for index, picture := range pictures {
			fmt.Printf("\n(%d/%d) ", index+1, total)
			pictureURL, err := url.Parse(picture)
			if err != nil {
				fmt.Printf("Error parsing picture %d (%s): %v", index, picture, err)
				continue
			}
			if !pictureURL.IsAbs() {
				if configuration.Base == "" {
					fmt.Print("Error: cannot load picture: its URL is relative and no -base flag was given")
					continue
				}
				pictureURL = joinURL(baseURL, pictureURL)
			}
			// fmt.Printf("Loading %s...", pictureURL.String())
			pictureName := path.Base(pictureURL.Path)
			if pictureName == "" || pictureName == "/" {
				fmt.Printf("Error: cannot determine picture name from path '%s'", pictureURL.Path)
				continue
			}
			output := uniqueName(path.Join(configuration.Output, pictureName))
			err = downloadPicture(pictureURL.String(), baseURL.String(), output, configuration.Agent, configuration.User, configuration.Password)
			if err != nil {
				fmt.Printf(" failed: %v", err)
			}
		}
		fmt.Println("")
	}
}

// uniqueName checks the file already exists: if yes it adds a (n) at the end
func uniqueName(filename string) string {
	if _, err := os.Stat(filename); err == nil || os.IsExist(err) {
		extension := path.Ext(filename)
		base := strings.TrimSuffix(filename, extension)
		index := 1
		for {
			filename = fmt.Sprintf("%s(%d)%s", base, index, extension)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				return filename
			}
			index++
		}
	}
	return filename
}
