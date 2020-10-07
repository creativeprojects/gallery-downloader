package main

import (
	"compress/gzip"
	"flag"
	"gallery-downloader/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func main() {
	var err error
	var configFile, root, port string
	flag.StringVar(&configFile, "config", "config.json", "client configuration file, used to compare http headers")
	flag.StringVar(&root, "root", "", "root where to serve your files from")
	flag.StringVar(&port, "port", "3000", "TCP port")
	flag.Parse()

	root = path.Clean(root)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", root)
	}

	cfg, err := config.LoadFileConfiguration(configFile)
	if err != nil {
		log.Fatalf("client configuration file not found")
	}

	log.Printf("Serving files from '%s' (use -root option to change the default)", root)
	fs := http.FileServer(http.Dir(root))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s - %s", r.Method, r.RequestURI, r.Referer())
		checkRequest(cfg.Browser, r)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fs.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fs.ServeHTTP(gzw, r)
	})

	log.Printf("Listening on :%s (use -port to change the default)...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func checkRequest(cfg config.Browser, request *http.Request) {
	// check user-agent
	if cfg.UserAgent != request.UserAgent() {
		log.Printf("WARNING: change user-agent in configuration for: '%s'", request.UserAgent())
	}

	// then check http headers
	if strings.HasSuffix(request.URL.Path, "html") {
		checkHeader(cfg.HTML, request.Header)
		return
	}
	if strings.HasSuffix(request.URL.Path, "jpg") {
		checkHeader(cfg.Picture, request.Header)
		return
	}
}

func checkHeader(cfg config.Element, headers http.Header) {
	// headers.Write(os.Stdout)
	// fmt.Println("")
}
