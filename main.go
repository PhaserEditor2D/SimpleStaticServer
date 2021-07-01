package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
)

var dir = "."
var port = 1999

func main() {

	flag.IntVar(&port, "port", port, "HTTP port")
	flag.StringVar(&dir, "dir", dir, "Content directory")
	flag.Parse()

	var err error

	if dir, err = filepath.Abs(dir); err != nil {

		log.Printf("%s", err.Error())
		os.Exit(1)
	}

	fileServer := http.FileServer(http.Dir(dir))
	cacheServer := cacheHandler(fileServer)
	zipServer := gziphandler.GzipHandler(cacheServer)

	http.Handle("/", zipServer)

	log.Printf("Serving directory '%s'", dir)
	log.Printf("Listening in http://127.0.0.1:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func cacheHandler(handler http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		isRoot := r.URL.Path == "/"

		if isRoot {

			w.Header().Set("Cache-Control", "no-store")

		} else {

			w.Header().Set("Cache-Control", "public, max-age=31104000, immutable")
		}

		handler.ServeHTTP(w, r)
	}
}
