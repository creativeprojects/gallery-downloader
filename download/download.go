package download

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"gallery-downloader/config"
	"gallery-downloader/headers"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// http client is global to the package, and instanciated on first use
var (
	transport *http.Transport
	client    *http.Client
)

// Config contains the configuration to download http files
type Config struct {
	Browser       config.Browser
	BaseURL       *url.URL
	Referer       string
	User          string
	Password      string
	Output        string
	WaitMin       int
	WaitMax       int
	Parallell     int
	SkipVerifyTLS bool
}

// Context contains the context to download http files
type Context struct {
	cfg Config
}

// NewContext creates a new Context with an http client.
// Any subsequent call to NewContext keeps the same http.Client already created
func NewContext(cfg Config) *Context {
	if transport == nil || client == nil {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.SkipVerifyTLS,
		}
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		}
		client = &http.Client{
			Transport: transport,
		}
	} else {
		transport.TLSClientConfig.InsecureSkipVerify = cfg.SkipVerifyTLS
	}

	return &Context{
		cfg: cfg,
	}
}

// HTML downloads an HTML page
func (c *Context) HTML(link string) ([]byte, error) {
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}
	c.setHTMLDownloadHeaders(request)

	response, err := client.Do(request)
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

// Pictures downloads a list of pictures
func (c *Context) Pictures(pictures []string) error {
	total := len(pictures)
	for index, picture := range pictures {
		fmt.Printf("\n(%d/%d) ", index+1, total)
		pictureURL, err := url.Parse(picture)
		if err != nil {
			fmt.Printf("Error parsing picture %d (%s): %v", index, picture, err)
			continue
		}
		if !pictureURL.IsAbs() {
			if c.cfg.BaseURL == nil || c.cfg.BaseURL.String() == "" {
				fmt.Print("Error: cannot load picture: its URL is relative and no -base flag was given")
				continue
			}
			pictureURL = joinURL(c.cfg.BaseURL, pictureURL)
		}
		pictureName := path.Base(pictureURL.Path)
		if pictureName == "" || pictureName == "/" {
			fmt.Printf("Error: cannot determine picture name from path '%s'", pictureURL.Path)
			continue
		}
		fmt.Printf("Loading %s...", pictureURL.String())
		output := uniqueName(path.Join(c.cfg.Output, pictureName))
		size, err := c.picture(pictureURL.String(), output)
		if err != nil {
			fmt.Printf(" failed: %v", err)
		} else {
			fmt.Printf(" loaded %d bytes", size)
			if size == 0 {
				// no need to keep an empty file
				fmt.Printf(" (not saved)")
				_ = os.Remove(output)
			}
			if c.cfg.WaitMax > 0 && c.cfg.WaitMax > c.cfg.WaitMin {
				wait := rand.Intn(c.cfg.WaitMax - c.cfg.WaitMin)
				fmt.Printf(" and wait %dms", wait+c.cfg.WaitMin)
				time.Sleep(time.Duration(wait+c.cfg.WaitMin) * time.Millisecond)
			}
		}
	}
	fmt.Println("")

	return nil
}

func (c *Context) picture(picture, output string) (int64, error) {
	request, err := http.NewRequest("GET", picture, nil)
	if err != nil {
		return 0, err
	}
	c.setPictureDownloadHeaders(request)

	// Output file
	if output == "" {
		output = os.DevNull
	}

	response, err := client.Do(request)
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

		// images shouldn't come back gzip encoded, but we never know
		reader := response.Body
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(response.Body)
			if err != nil {
				return 0, err
			}
		}

		size, err := io.Copy(outputFile, reader)
		if err != nil {
			return size, err
		}
		return size, nil
	}
	return 0, fmt.Errorf("HTTP %s", response.Status)
}

func (c *Context) setHTMLDownloadHeaders(request *http.Request) {
	c.setCommonHeaders(request)
	for name, value := range c.cfg.Browser.HTML.Headers {
		request.Header.Set(name, value)
	}
}

func (c *Context) setPictureDownloadHeaders(request *http.Request) {
	c.setCommonHeaders(request)
	for name, value := range c.cfg.Browser.Picture.Headers {
		request.Header.Set(name, value)
	}
}

func (c *Context) setCommonHeaders(request *http.Request) {
	for name, value := range c.cfg.Browser.Default.Headers {
		request.Header.Set(name, value)
	}

	if request.URL.Scheme == "http" {
		for name, value := range c.cfg.Browser.HTTP.Headers {
			request.Header.Set(name, value)
		}
	} else if request.URL.Scheme == "https" {
		for name, value := range c.cfg.Browser.HTTPS.Headers {
			request.Header.Set(name, value)
		}
	}

	if c.cfg.Referer != "" {
		request.Header.Set(headers.Referer, c.cfg.Referer)
	}

	if c.cfg.User != "" && c.cfg.Password != "" {
		request.SetBasicAuth(c.cfg.User, c.cfg.Password)
	}
}
