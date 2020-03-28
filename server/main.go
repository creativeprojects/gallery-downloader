package main

import (
	"flag"
	"log"
	"net/http"
	"path"
)

func main() {
	var root, port string
	flag.StringVar(&root, "root", "", "root where to serve your files from")
	flag.StringVar(&port, "port", "3000", "TCP port")
	flag.Parse()

	root = path.Clean(root)

	log.Printf("Serving files from '%s' (use -root option to change the default)", root)
	fs := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		fs.ServeHTTP(w, r)
	})

	log.Printf("Listening on :%s (use -port to change the default)...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
