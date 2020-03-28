package main

import (
	"compress/gzip"
	"flag"
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
	var root, port string
	flag.StringVar(&root, "root", "", "root where to serve your files from")
	flag.StringVar(&port, "port", "3000", "TCP port")
	flag.Parse()

	root = path.Clean(root)

	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Fatalf("'%s' does not exist", root)
	}

	log.Printf("Serving files from '%s' (use -root option to change the default)", root)
	fs := http.FileServer(http.Dir(root))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %d - %s - %s", r.Method, len(r.Header), r.RequestURI, r.Referer())
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
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
