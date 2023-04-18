package main

import (
	"context"
	"errors"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"
)

var baseDirectory string

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join(baseDirectory, r.URL.Path)
	file, err := os.Stat(filePath)
	if err == nil && !file.IsDir() {
		// file exists
		gz := filePath + ".gz"
		file, err := os.Stat(gz)
		t := mime.TypeByExtension(filepath.Ext(filePath))
		if err == nil && !file.IsDir() && t != "" {
			w.Header().Add("Content-Encoding", "gzip")
			w.Header().Add("Content-Type", t)
			http.ServeFile(w, r, gz)
		} else {
			http.ServeFile(w, r, filePath)
		}
	} else {
		// file does not exist
		index := path.Join(baseDirectory, "index.html")
		file, err := os.Stat(index)
		if err == nil && !file.IsDir() {
			// index.html exists
			http.ServeFile(w, r, index)
		} else {
			// index.html does not exist
			http.NotFound(w, r)
		}
	}
}

func main() {
	listeningPort, ok := os.LookupEnv("PORT")
	if !ok {
		listeningPort = "5050"
	}

	listeningHost, ok := os.LookupEnv("HOST")
	if !ok {
		listeningHost = "127.0.0.1"
	}

	baseDirectory, ok = os.LookupEnv("BASE_DIRECTORY")
	if !ok {
		baseDirectory = "."
	}

	router := http.NewServeMux()
	router.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    net.JoinHostPort(listeningHost, listeningPort),
		Handler: router,
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal

		log.Println("Received shutdown signal, shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("Error occured during shutting down HTTP server: %s", err.Error())
		}
	}()

	log.Printf("HTTP server listening on %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error occured during listening to HTTP server: %s", err.Error())
	}
}
