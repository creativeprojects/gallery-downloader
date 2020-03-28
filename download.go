package main

import (
	"io"
	"net/http"
	"os"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{}
}

func downloadPicture(picture, referer, output, agent, user, password string) error {
	request, err := http.NewRequest("GET", picture, nil)
	if err != nil {
		return err
	}
	// These should all move to a configuration file!
	request.Header.Set("Accept", "image/webp,*/*")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("DNT", "1")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Referer", referer)
	request.Header.Set("User-Agent", agent)
	//

	if user != "" && password != "" {
		request.SetBasicAuth(user, password)
	}

	// Output file
	if output == "" {
		output = os.DevNull
	}
	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	_, err = io.Copy(outputFile, response.Body)
	response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
