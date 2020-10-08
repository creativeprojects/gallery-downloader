package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"gallery-downloader/config"
	"gallery-downloader/headers"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	var configFile, root, httpPort, httpsPort, certFile, keyFile string
	var verbose bool
	flag.StringVar(&configFile, "config", "config.json", "client configuration file, used to compare http headers")
	flag.StringVar(&root, "root", "", "root where to serve your files from")
	flag.StringVar(&httpPort, "http", "3000", "TCP port for HTTP requests")
	flag.StringVar(&httpsPort, "https", "3001", "TCP port for HTTPS requests")
	flag.StringVar(&certFile, "cert", "cert/localhost.cert.pem", "certificate to serve HTTPS requests")
	flag.StringVar(&keyFile, "key", "cert/localhost.key.pem", "private key to serve HTTPS requests")
	flag.BoolVar(&verbose, "v", false, "display debugging information (mostly HTTP request headers)")
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

	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s - %s - %s", r.Proto, r.Method, r.RequestURI, r.Referer())
		checkRequest(cfg.Browser, r, verbose)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fs.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fs.ServeHTTP(gzw, r)
	}

	http.HandleFunc("/", handleRequest)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	log.Printf("HTTP: listening on :%s (use -http to change the default port)...", httpPort)
	go func() {
		err = http.ListenAndServe(":"+httpPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if certFile != "" && keyFile != "" && httpsPort != "" {
		log.Printf("HTTPS: listening on :%s (use -https to change the default port)...", httpsPort)
		go func() {
			err = http.ListenAndServeTLS(":"+httpsPort, certFile, keyFile, nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	<-stop
	fmt.Println("")
}

func checkRequest(cfg config.Browser, request *http.Request, verbose bool) {
	// check user-agent
	if cfg.Default.Headers[headers.UserAgent] != request.UserAgent() {
		log.Printf("WARNING: change user-agent in configuration for: '%s'", request.UserAgent())
	}

	// then check http headers
	if strings.HasSuffix(request.URL.Path, "html") {
		checkHeader(cfg.HTML, request.Header, verbose)
		return
	}
	if strings.HasSuffix(request.URL.Path, "jpg") {
		checkHeader(cfg.Picture, request.Header, verbose)
		return
	}
}

func checkHeader(cfg config.Group, headers http.Header, verbose bool) {
	if !verbose {
		return
	}
	headers.Write(os.Stdout)
	fmt.Println("")
}
