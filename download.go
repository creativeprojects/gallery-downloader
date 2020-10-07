package main

import (
	"compress/gzip"
	"fmt"
	"gallery-downloader/config"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// DownloadConfig contains the configuration to download http files
type DownloadConfig struct {
	Client   *http.Client
	BaseURL  *url.URL
	Browser  config.Browser
	Referer  string
	User     string
	Password string
	Output   string
	WaitMin  int
	WaitMax  int
}

// NewDownloadConfig creates a new DownloadConfig with an http client
func NewDownloadConfig(baseURL *url.URL, referer, user, password, output string, browser config.Browser, waitMin, waitMax int) DownloadConfig {
	return DownloadConfig{
		Client:   &http.Client{},
		BaseURL:  baseURL,
		Referer:  referer,
		Browser:  browser,
		User:     user,
		Password: password,
		Output:   output,
		WaitMin:  waitMin,
		WaitMax:  waitMax,
	}
}

func downloadHTML(link string, downloadConfig DownloadConfig) ([]byte, error) {
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	setHTMLDownloadHeaders(request, downloadConfig)

	response, err := downloadConfig.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 400 {
		reader := response.Body
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(response.Body)
			if err != nil {
				return nil, err
			}
		}
		buffer, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return buffer, nil
	}
	return nil, fmt.Errorf("HTTP %s", response.Status)
}

func downloadPictures(pictures []string, downloadConfig DownloadConfig) error {
	total := len(pictures)
	for index, picture := range pictures {
		fmt.Printf("\n(%d/%d) ", index+1, total)
		pictureURL, err := url.Parse(picture)
		if err != nil {
			fmt.Printf("Error parsing picture %d (%s): %v", index, picture, err)
			continue
		}
		if !pictureURL.IsAbs() {
			if downloadConfig.BaseURL == nil || downloadConfig.BaseURL.String() == "" {
				fmt.Print("Error: cannot load picture: its URL is relative and no -base flag was given")
				continue
			}
			pictureURL = joinURL(downloadConfig.BaseURL, pictureURL)
		}
		pictureName := path.Base(pictureURL.Path)
		if pictureName == "" || pictureName == "/" {
			fmt.Printf("Error: cannot determine picture name from path '%s'", pictureURL.Path)
			continue
		}
		fmt.Printf("Loading %s...", pictureURL.String())
		output := uniqueName(path.Join(downloadConfig.Output, pictureName))
		size, err := downloadPicture(pictureURL.String(), output, downloadConfig)
		if err != nil {
			fmt.Printf(" failed: %v", err)
		} else {
			fmt.Printf(" loaded %d bytes", size)
			if downloadConfig.WaitMax > 0 && downloadConfig.WaitMax > downloadConfig.WaitMin {
				wait := rand.Intn(downloadConfig.WaitMax - downloadConfig.WaitMin)
				fmt.Printf(" and wait %dms", wait+downloadConfig.WaitMin)
				time.Sleep(time.Duration(wait+downloadConfig.WaitMin) * time.Millisecond)
			}
		}
	}
	fmt.Println("")

	return nil
}

func downloadPicture(picture, output string, downloadConfig DownloadConfig) (int64, error) {
	request, err := http.NewRequest("GET", picture, nil)
	if err != nil {
		return 0, err
	}
	setPictureDownloadHeaders(request, downloadConfig)

	// Output file
	if output == "" {
		output = os.DevNull
	}

	response, err := downloadConfig.Client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 400 {
		outputFile, err := os.Create(output)
		if err != nil {
			return 0, err
		}
		defer outputFile.Close()

		size, err := io.Copy(outputFile, response.Body)
		if err != nil {
			return size, err
		}
		return size, nil
	}
	return 0, fmt.Errorf("HTTP %s", response.Status)
}

func setHTMLDownloadHeaders(request *http.Request, downloadConfig DownloadConfig) {
	for name, value := range downloadConfig.Browser.HTML.Headers {
		request.Header.Set(name, value)
	}
	setCommonHeaders(request, downloadConfig)
}

func setPictureDownloadHeaders(request *http.Request, downloadConfig DownloadConfig) {
	for name, value := range downloadConfig.Browser.Picture.Headers {
		request.Header.Set(name, value)
	}
	setCommonHeaders(request, downloadConfig)
}

func setCommonHeaders(request *http.Request, downloadConfig DownloadConfig) {
	if downloadConfig.Referer != "" {
		request.Header.Set("Referer", downloadConfig.Referer)
	}
	request.Header.Set("User-Agent", downloadConfig.Browser.UserAgent)

	if downloadConfig.User != "" && downloadConfig.Password != "" {
		request.SetBasicAuth(downloadConfig.User, downloadConfig.Password)
	}
}
