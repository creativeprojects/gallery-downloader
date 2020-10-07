package main

import (
	"bytes"
	"flag"
	"gallery-downloader/config"
	"gallery-downloader/scan"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	var err error

	setLogger()
	flags := loadFlags()

	checkSource(flags)
	checkType(flags)
	checkOutput(flags)

	cfg, err := config.LoadFileConfiguration(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Cannot load configuration: %v", err)
	}

	var baseURL = &url.URL{}
	if flags.Base != "" {
		baseURL, err = url.Parse(flags.Base)
		if err != nil {
			log.Fatal("Error: -base value is not a parsable URL")
		}
		if !baseURL.IsAbs() {
			log.Fatal("Error: -base value is not an absolute URL")
		}
	}

	sourceURL, err := url.Parse(flags.Source)
	if err != nil {
		log.Fatalf("Error parsing source URL: %v", err)
	}

	if sourceURL.Scheme == "" {
		downloadPicturesFromLocalGalleryFile(flags.Source, baseURL, flags, cfg.Browser)
	} else {
		downloadPicturesFromRemoteGallery(sourceURL, flags, cfg.Browser)
	}

}

func setLogger() {
	log.SetFlags(0)
}

func checkSource(flags Flags) {
	if flags.Source == "" {
		flag.Usage()
		log.Fatal("\nError: missing HTML source (-source)")
	}
}

func checkType(flags Flags) {
	if flags.Type == "" {
		flag.Usage()
		log.Fatal("\nError: missing gallery type (-type)")
	}
	for _, galleryType := range scan.AvailableGalleryScanners {
		if galleryType == flags.Type {
			// nothing else to check
			return
		}
	}
	log.Fatalf("\nError: unknown gallery type. Known types are: %s", strings.Join(scan.AvailableGalleryScanners[:], ", "))
}

func checkOutput(flags Flags) {
	if flags.Output == "" {
		flag.Usage()
		log.Fatal("\nError: missing output folder (-output)")
	}
	if stat, err := os.Stat(flags.Output); err == nil || os.IsExist(err) {
		if !stat.IsDir() {
			log.Fatalf("Output '%s' exists but is not a directory", flags.Output)
		}
	}
	if _, err := os.Stat(flags.Output); os.IsNotExist(err) {
		err = os.MkdirAll(flags.Output, 0755)
		if err != nil {
			log.Fatalf("Cannot create output directory: %v", err)
		}
	}
}

func downloadPicturesFromLocalGalleryFile(sourceFile string, baseURL *url.URL, flags Flags, browserConfig config.Browser) {
	// Let's consider this is a file on disk
	sourcefile, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("Error: cannot open HTML source file: %v", err)
	}
	defer sourcefile.Close()
	pictures := scanPictures(sourcefile, flags)
	if pictures == nil || len(pictures) == 0 {
		log.Println("No picture found in the HTML source!")
	}

	downloadConfig := NewDownloadConfig(baseURL, flags.Referer, flags.User, flags.Password, flags.Output, browserConfig, flags.WaitMin, flags.WaitMax)
	downloadPictures(pictures, downloadConfig)
}

func downloadPicturesFromRemoteGallery(sourceURL *url.URL, flags Flags, browserConfig config.Browser) {
	// We need to download the remote HTML file
	downloadConfig := NewDownloadConfig(nil, flags.Referer, flags.User, flags.Password, "", browserConfig, 0, 0)
	buffer, err := downloadHTML(flags.Source, downloadConfig)
	if err != nil {
		log.Fatalf("Error: cannot download HTML source file: %v", err)
	}
	pictures := scanPictures(bytes.NewReader(buffer), flags)
	if pictures == nil || len(pictures) == 0 {
		ioutil.WriteFile(path.Join(flags.Output, "index.html"), buffer, 0644)
		log.Println("No picture found in the HTML source. HTML file saved.")
	}

	downloadConfig = NewDownloadConfig(sourceURL, flags.Source, flags.User, flags.Password, flags.Output, browserConfig, flags.WaitMin, flags.WaitMax)
	downloadPictures(pictures, downloadConfig)
}

func scanPictures(source io.ReadSeeker, flags Flags) []string {
	var err error
	var pictures []string
	log.Printf("using gallery scanner: %s", flags.Type)
	for _, scanner := range scan.GalleryScanners[flags.Type] {
		source.Seek(0, 0)
		pictures, err = scanner(source)
		if err != nil {
			log.Fatalf("Error: cannot parse HTML source file: %v", err)
		}
		if len(pictures) > 1 {
			// no need to try another one
			return pictures
		}
	}
	return pictures
}
