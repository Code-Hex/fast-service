package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Code-Hex/fast-service/internal/config"
	"github.com/Code-Hex/fast-service/internal/randomer"
)

const maxSize = 26214400 // 25MB

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Read configurations from environmental variables.
	env, err := config.ReadFromEnv()
	if err != nil {
		log.Printf("failed to read environment variables: %s", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/download", downloadHandler())
	mux.Handle("/upload", uploadHandler())

	srv := http.Server{
		Handler: mux,
	}

	addr := fmt.Sprintf(":%d", env.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen port: %s", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("serve err: %s", err)
			return
		}
	}()

	// waiting for SIGTERM or Interrupt signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		log.Printf("received SIGTERM, exiting server gracefully")
	case <-ctx.Done():
	}
	log.Printf("shutdown servers")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown HTTP server: %s", err)
		return
	}
}

func downloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		queries := r.URL.Query()
		size := queries.Get("size")
		max, err := strconv.Atoi(size)
		if err != nil {
			max = maxSize
		}
		if _, err := io.CopyN(w, randomer.New(), int64(max)); err != nil {
			log.Printf("failed to write random file: %s", err)
			return
		}
	}
}

func uploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			log.Printf("invalid content type: %s", contentType)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.ContentLength == 0 {
			log.Printf("invalid content length")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentLength := r.ContentLength
		if contentLength > maxSize {
			contentLength = maxSize
		}
		if _, err := io.CopyN(ioutil.Discard, r.Body, contentLength); err != nil {
			log.Printf("failed to write body: %s", err)
			return
		}
	}
}
