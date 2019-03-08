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

	"go.uber.org/zap"

	"github.com/Code-Hex/fast-service/internal/config"
	"github.com/Code-Hex/fast-service/internal/logger"
	"github.com/Code-Hex/fast-service/internal/server"
)

const maxSize = 26214400 // 25MB

func main() {
	// Read configurations from environmental variables.
	env, err := config.ReadFromEnv()
	if err != nil {
		log.Fatalf("failed to read environment variables: %s", err)
	}

	// Setup new zap logger. This logger should be used for all logging in this service.
	// The log level can be updated via environment variables.
	l, err := logger.New(env.LogLevel)
	if err != nil {
		log.Fatalf("failed to prepare logger: %s", err)
	}

	if err := _main(env, l); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func _main(env *config.Env, l *zap.Logger) error {
	mux := server.NewMux(l)
	mux.Handle("/download", downloadHandler())
	mux.Handle("/upload", uploadHandler())

	srv := server.New(mux)

	addr := fmt.Sprintf(":%d", env.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen port: %s", err)
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
		return fmt.Errorf("failed to gracefully shutdown HTTP server: %s", err)
	}
	return nil
}

func downloadHandler() http.HandlerFunc {
	src := rand.NewSource(0)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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
		if _, err := io.CopyN(w, rand.New(src), int64(max)); err != nil {
			logger.Error(ctx, "failed to write random file: %s", zap.Error(err))
			return
		}
	}
}

func uploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			logger.Warn(ctx, "invalid content type", zap.String("Content-Type", contentType))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		contentLength := r.ContentLength
		if contentLength > maxSize {
			contentLength = maxSize
		}
		if _, err := io.CopyN(ioutil.Discard, r.Body, contentLength); err != nil {
			logger.Warn(ctx, "failed to write body", zap.Error(err))
			return
		}
	}
}
