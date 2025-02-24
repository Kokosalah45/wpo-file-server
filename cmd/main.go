package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

type WeatherConfig struct {
	Rain int `json:"rain"`
	Snow int `json:"snow"`
	Sun int `json:"sun"`
	Heat int `json:"heat"`
	Wind int `json:"wind"`
	
}

func Start() {
	mux.HandleFunc("/config/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Cache-Control", "max-age=120, must-revalidate")
		w.Header().Add("Content-Type", "application/json")
		config := WeatherConfig{
			Rain: 10,
			Snow: 20,
			Sun: 30,
			Heat: 40,
			Wind: 50,
		}
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(config)
	})

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {

		filePath := r.URL.Path[1:]

		fmt.Println(filePath)
		file, err := os.ReadFile(filePath)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		etag := r.Header.Get("If-None-Match")

		fileMd5 := fmt.Sprintf("%x", md5.Sum(file))
		fileExt := filepath.Ext(filePath)
		cacheControlConfig := fileCacheStrategyConfig[fileExt]


		if etag == fileMd5 {
			w.Header().Add("Cache-Control", cacheControlConfig)
			w.Header().Add("Etag", fileMd5)
			w.Header().Add("Content-Length", fmt.Sprintf("%d", len(file)))
			w.Header().Add("Content-Type" ,  mimeTypes[fileExt])
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
