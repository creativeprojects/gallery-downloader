package download

import (
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"gallery-downloader/config"
	"gallery-downloader/headers"
	"io"
	"io/ioutil"
	"log"
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

type job struct {
	picture string
	index   int
	total   int
}

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
	Parallel      int
	SkipVerifyTLS bool
	Progress      func(Progress)
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
func (c *Context) Pictures(pictures []string) {
	total := len(pictures)
	if c.cfg.Parallel < 2 {
		// simple case of synchronous download
		for index, picture := range pictures {
			c.picture(picture, index, total)
		}
		return
	}

	jobs := make(chan job, total)
	results := make(chan interface{}, total)

	for w := 1; w <= c.cfg.Parallel; w++ {
		go c.pictureWorker(w, jobs, results)
	}

	for index, picture := range pictures {
		jobs <- job{picture, index, total}
	}
	close(jobs)

	// get all results (could also do with a waitgroup)
	for a := 1; a <= total; a++ {
		<-results
	}
}

func (c *Context) pictureWorker(id int, jobs <-chan job, results chan<- interface{}) {
	log.Printf("Creating worker %d", id)
	for j := range jobs {
		c.picture(j.picture, j.index, j.total)
		results <- nil
	}
	log.Printf("Worker %d finished", id)
}

func (c *Context) picture(picture string, index, total int) {
	pictureURL, err := url.Parse(picture)
	if err != nil {
		if c.cfg.Progress != nil {
			c.cfg.Progress(Progress{
				FileID:     index,
				TotalFiles: total,
				Event:      EventError,
				Err:        fmt.Errorf("invalid picture URL: %w", err),
			})
		}
		return
	}
	if !pictureURL.IsAbs() {
		if c.cfg.BaseURL == nil || c.cfg.BaseURL.String() == "" {
			if c.cfg.Progress != nil {
				c.cfg.Progress(Progress{
					FileID:     index,
					TotalFiles: total,
					URL:        pictureURL.String(),
					Event:      EventError,
					Err:        errors.New("cannot load picture: its URL is relative and no -base flag was given"),
				})
			}
			return
		}
		pictureURL = joinURL(c.cfg.BaseURL, pictureURL)
	}
	pictureName := path.Base(pictureURL.Path)
	if pictureName == "" || pictureName == "/" {
		if c.cfg.Progress != nil {
			c.cfg.Progress(Progress{
				FileID:     index,
				TotalFiles: total,
				URL:        pictureURL.String(),
				Event:      EventError,
				Err:        fmt.Errorf("cannot determine picture name from path '%s'", pictureURL.Path),
			})
		}
		return
	}
	if c.cfg.Progress != nil {
		c.cfg.Progress(Progress{
			FileID:     index,
			TotalFiles: total,
			URL:        pictureURL.String(),
			Event:      EventStart,
		})
	}
	output := uniqueName(path.Join(c.cfg.Output, pictureName))
	size, err := c.downloadPicture(pictureURL.String(), output)
	if err != nil {
		if c.cfg.Progress != nil {
			c.cfg.Progress(Progress{
				FileID:     index,
				TotalFiles: total,
				URL:        pictureURL.String(),
				Event:      EventError,
				Err:        err,
			})
		}
	} else {
		progress := Progress{
			FileID:     index,
			TotalFiles: total,
			URL:        pictureURL.String(),
			Event:      EventFinished,
			Downloaded: size,
		}
		if size == 0 {
			// no need to keep an empty file
			progress.Event = EventNotSaving
			_ = os.Remove(output)
		}
		if c.cfg.WaitMax > 0 && c.cfg.WaitMax > c.cfg.WaitMin {
			wait := rand.Intn(c.cfg.WaitMax - c.cfg.WaitMin)
			progress.Wait = wait + c.cfg.WaitMin
			// now send the complete progress report
			if c.cfg.Progress != nil {
				c.cfg.Progress(progress)
			}
			// and wait
			time.Sleep(time.Duration(wait+c.cfg.WaitMin) * time.Millisecond)
		} else {
			// now send the complete progress report
			if c.cfg.Progress != nil {
				c.cfg.Progress(progress)
			}
		}
	}
}

func (c *Context) downloadPicture(picture, output string) (int64, error) {
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
