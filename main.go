package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
)

const ver = "1.0.0"

var dir = "."
var port = 1999

func main() {

	fmt.Printf("Simple Static Server v%s\n\n", ver)

	flag.Usage = func() {

		fmt.Fprintf(flag.CommandLine.Output(), "\nCommon usage:\n\n")
		fmt.Println("$ SimpleStaticServer")
		fmt.Println("$ SimpleStaticServer -dir /path/to/static/content -port 80")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println()

		flag.PrintDefaults()

		fmt.Println()
		fmt.Println("Contact: Support Team <developers@phasereditor2d.com>")
		fmt.Println("Copyrights (c) Arian Fornaris <arian@phasereditor2d.com>")
		fmt.Println()
	}

	flag.IntVar(&port, "port", port, "HTTP port")
	flag.StringVar(&dir, "dir", dir, "Content directory")
	flag.Parse()

	var err error

	if dir, err = filepath.Abs(dir); err != nil {

		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	fileServer := http.FileServer(NewCustomFileSystem(dir))
	cacheServer := cacheHandler(fileServer)
	zipServer := gziphandler.GzipHandler(cacheServer)

	http.Handle("/", zipServer)

	fmt.Printf("Serving directory '%s'\n", dir)
	fmt.Printf("Listening http://127.0.0.1:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// customFileSystem a file system for handling "file not found".
type customFileSystem struct {
	baseFileSystem http.FileSystem
}

// Open the requested file. If the file is not found, it opens the root file (`/index.html`).
func (fs customFileSystem) Open(path string) (http.File, error) {

	file, err := fs.baseFileSystem.Open(path)

	if err != nil {

		// file not found, let's serve the root file

		return fs.baseFileSystem.Open("/index.html")
	}

	return file, nil
}

// NewCustomFileSystem creates a new custom file system.
func NewCustomFileSystem(path string) *customFileSystem {

	return &customFileSystem{
		baseFileSystem: http.Dir(path),
	}
}

// cacheHandler set `Cache-Control` header.
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
