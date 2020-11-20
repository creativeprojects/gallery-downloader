package main

import (
	"flag"
	"fmt"
	"gallery-downloader/config"
	"gallery-downloader/download"
	"gallery-downloader/scan"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/andybalholm/cascadia"
)

func main() {
	var err error

	setLogger()
	flags := loadFlags()

	checkSource(flags)
	checkType(flags)
	checkOutput(flags)
	checkParallel(flags)

	cfg, err := config.LoadFileConfiguration(flags.ConfigFile)
	if err != nil {
		log.Fatalf("Error: cannot load configuration: %v", err)
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
		downloadPicturesFromLocalGalleryFile(flags.Source, baseURL, flags, cfg)
	} else {
		downloadPicturesFromRemoteGallery(sourceURL, flags, cfg)
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
			log.Fatalf("\nError: output '%s' exists but is not a directory", flags.Output)
		}
	}
	if _, err := os.Stat(flags.Output); os.IsNotExist(err) {
		err = os.MkdirAll(flags.Output, 0755)
		if err != nil {
			log.Fatalf("\nError: cannot create output directory: %v", err)
		}
	}
}

func checkParallel(flags Flags) {
	if flags.Parallel > 1 &&
		(flags.WaitMin > 0 || flags.WaitMax > 0) {
		log.Fatal("\nError: cannot use parallel download with -min-wait and/or -max-wait parameters at the same time")
	}
}

func downloadPicturesFromLocalGalleryFile(sourceFile string, baseURL *url.URL, flags Flags, cfg *config.Configuration) {
	// Let's consider this is a file on disk
	sourcefile, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("Error: cannot open HTML source file: %v", err)
	}
	defer sourcefile.Close()
	buffer, err := ioutil.ReadAll(sourcefile)
	if err != nil {
		log.Fatalf("Error cannot read gallery file: %s", err)
	}
	pictures := scanImages(buffer, flags, cfg)
	if pictures == nil || len(pictures) == 0 {
		log.Println("No picture found in the HTML source!")
	}

	downloadContext := download.NewContext(download.Config{
		BaseURL:       baseURL,
		Referer:       flags.Referer,
		User:          flags.User,
		Password:      flags.Password,
		Output:        flags.Output,
		Browser:       cfg.Browser,
		WaitMin:       flags.WaitMin,
		WaitMax:       flags.WaitMax,
		SkipVerifyTLS: flags.InsecureTLS,
		Parallel:      flags.Parallel,
		Progress:      handleProgress,
	})
	downloadContext.Pictures(pictures)
}

func downloadPicturesFromRemoteGallery(sourceURL *url.URL, flags Flags, cfg *config.Configuration) {
	// We need to download the remote HTML file
	downloadContext := download.NewContext(download.Config{
		Referer:       flags.Referer,
		User:          flags.User,
		Password:      flags.Password,
		Browser:       cfg.Browser,
		SkipVerifyTLS: flags.InsecureTLS,
		Progress:      handleProgress,
	})
	buffer, err := downloadContext.HTML(flags.Source)
	if err != nil {
		log.Fatalf("Error: cannot download HTML source file: %v", err)
	}
	pictures := scanImages(buffer, flags, cfg)
	if pictures == nil || len(pictures) == 0 {
		ioutil.WriteFile(path.Join(flags.Output, "index.html"), buffer, 0644)
		log.Println("No picture found in the HTML source. HTML file saved as index.html")
	}

	downloadContext = download.NewContext(download.Config{
		User:          flags.User,
		Password:      flags.Password,
		Browser:       cfg.Browser,
		SkipVerifyTLS: flags.InsecureTLS,
		BaseURL:       sourceURL,
		Referer:       flags.Source,
		Output:        flags.Output,
		WaitMin:       flags.WaitMin,
		WaitMax:       flags.WaitMax,
		Parallel:      flags.Parallel,
		Progress:      handleProgress,
	})
	downloadContext.Pictures(pictures)
}

func handleProgress(progress download.Progress) {
	count := ""
	if progress.TotalFiles > 0 {
		count = fmt.Sprintf("(%d/%d) ", progress.FileID+1, progress.TotalFiles)
	}
	message := ""
	switch progress.Event {
	case download.EventStart:
		message = fmt.Sprintf("download starting: '%s'", progress.URL)
	case download.EventFinished:
		message = fmt.Sprintf("  finished downloading %d bytes", progress.Downloaded)
	case download.EventNotSaving:
		message = fmt.Sprintf("  not saving file of %d bytes", progress.Downloaded)
	case download.EventError:
		message = fmt.Sprintf("error: %s", progress.Err)
	}
	wait := ""
	if progress.Wait > 0 {
		wait = fmt.Sprintf(" and wait for %dms", progress.Wait)
	}
	log.Printf("%s%s%s", count, message, wait)
}

func scanImages(source []byte, flags Flags, cfg *config.Configuration) []string {
	var pictures []string
	log.Printf("Using gallery scanner: %s", flags.Type)

	if flags.Type == scan.AutoDetect || flags.Type == scan.ConfigProfiles {
		// first pass, use regexp profiles from configuration
		pictures, err := detectFromProfiles(cfg.Profiles, source)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if len(pictures) > 3 {
			return pictures
		}
	}

	// var err error
	// buffer := bytes.NewReader(source)
	// // use the legacy scanners
	// for _, scanner := range scan.GalleryScanners[flags.Type] {
	// 	pictures, err = scanner(buffer)
	// 	if err != nil {
	// 		log.Fatalf("Error: cannot parse HTML source file: %v", err)
	// 	}
	// 	if len(pictures) > 3 {
	// 		// no need to try another one
	// 		return pictures
	// 	}
	// 	// rewind source stream
	// 	_, err = buffer.Seek(0, io.SeekStart)
	// 	if err != nil {
	// 		log.Fatalf("Error: cannot rewind stream: %s", err)
	// 	}
	// }
	return pictures
}

func detectFromProfiles(profiles []config.Profile, source []byte) ([]string, error) {
	// current minimum priority to choose from
	priority := -1
	for {
		// on each turn, we need to pick the smallest priority number between "priority" and "next"
		// "next" represents the smallest number bigger than "priority"
		next := priority
		run := -1
		for i, profile := range profiles {
			// first time here: next == priority
			if profile.Priority > priority && (next == priority || profile.Priority < next) {
				next = profile.Priority
				run = i
			}
		}

		// no profile was chosen, it means we've been through them all
		if run < 0 {
			break
		}

		profile := profiles[run]

		generator, err := newMatcher(profile.DetectGenerator)
		if err != nil {
			return nil, fmt.Errorf("profile %s: cannot compile generator %s: %w", profile.Name, profile.DetectGenerator, err)
		}

		gallery, err := newMatcher(profile.DetectGallery)
		if err != nil {
			return nil, fmt.Errorf("profile %s: cannot compile generator %s: %w", profile.Name, profile.DetectGallery, err)
		}

		image, err := newMatcher(profile.DetectImage)
		if err != nil {
			return nil, fmt.Errorf("profile %s: cannot compile generator %s: %w", profile.Name, profile.DetectImage, err)
		}
		if image == nil {
			return nil, fmt.Errorf("profile %s: missing detectImage", profile.Name)
		}

		scanner := scan.NewGallery(scan.Config{
			Name:            profile.Name,
			DetectGenerator: generator,
			DetectGallery:   gallery,
			DetectImage:     image,
		}, source)

		if scanner.Match() {
			images := scanner.Found()
			if len(images) > 3 {
				log.Printf("Found %d images using profile %s (#%d)", len(images), profile.Name, run+1)
				generatedBy := scanner.GeneratedBy()
				if generatedBy != "" {
					log.Printf("Gallery generated by %s", generatedBy)
				}
				return images, nil
			}
		}
		// on next turn, don't pick anything with less priority than this one
		priority = profiles[run].Priority
	}
	return nil, nil
}

func newMatcher(cfg config.Parser) (scan.Matcher, error) {
	matcherType := strings.ToLower(cfg.Type)

	if strings.HasPrefix(matcherType, "regex") {
		pattern, err := regexp.Compile(cfg.Match)
		if err != nil {
			return nil, err
		}
		return scan.NewRegexpMatcher(pattern), nil
	}

	if strings.HasPrefix(matcherType, "sel") || strings.HasPrefix(matcherType, "css") {
		sel, err := cascadia.Parse(cfg.Match)
		if err != nil {
			return nil, err
		}
		return scan.NewSelectorMatcher(sel, cfg.Attribute), nil
	}
	return nil, nil
}
