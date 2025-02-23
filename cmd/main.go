package app

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var fileCacheStrategyConfig = map[string]string{
	".js":   "max-age=31536000, immutable",
	".css":  "max-age=31536000, immutable",
	".jpg":  "max-age=31536000, immutable",
	".jpeg": "max-age=31536000, immutable",
	".png":  "max-age=31536000, immutable",
	".gif":  "max-age=31536000, immutable",
	".svg":  "max-age=31536000, immutable",
	".ico":  "max-age=31536000, immutable",
	".html": "no-cache",
	".json": "private, no-store",
}

var mimeTypes = map[string]string{
	".js":   "text/javascript",
	".css":  "text/css",
	".jpg":  "image/jpeg",
	".html": "text/html",
}
var mux = http.NewServeMux()

func Start() {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		const prefix = "static"
		filePath := r.URL.Path[1:]

		fmt.Println("Request URL: ", r.URL.Path)
		fmt.Println("File Path: ", filePath)

		hasPrefix := strings.HasPrefix(filePath, prefix)

		if !hasPrefix {
			filePath = path.Join(prefix, r.URL.Path)
		}

		file, err := os.ReadFile(filePath)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		etag := r.Header.Get("If-None-Match")

		fileMd5 := fmt.Sprintf("%x", md5.Sum(file))
		cacheControlConfig := "no-cache"


		if etag == fileMd5 {
			w.Header().Add("Cache-Control", cacheControlConfig)
			w.Header().Add("Etag", fileMd5)
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(file)))
			w.Header().Add("Content-Type" ,  mimeTypes[filepath.Ext(filePath)])
			w.WriteHeader(http.StatusNotModified)

			return 
		}

		w.Header().Add("Cache-Control", cacheControlConfig)
		w.Header().Add("Etag", fileMd5)
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(file)))
		w.Header().Add("Content-Type" , mimeTypes[filepath.Ext(filePath)])

		w.Write(file)
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	} else {
		fmt.Println("Server is running on port 8080")
	}

}
