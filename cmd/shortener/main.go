package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/N0rkton/shortener/internal/app/config"
	"github.com/N0rkton/shortener/internal/app/handlers"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

const workerCount = 10

func main() {
	var wg sync.WaitGroup
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", handlers.JSONIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten/batch", handlers.Batch).Methods(http.MethodPost)
	router.HandleFunc("/ping", handlers.PingDB).Methods(http.MethodGet)
	router.HandleFunc("/{id}", handlers.RedirectTo).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", handlers.ListURL).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", handlers.DeleteURL).Methods(http.MethodDelete)
	var srv = http.Server{Addr: config.GetServerAddress(), Handler: handlers.GzipHandle(router)}

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	go func() {
		defer func() {
			wg.Wait()
		}()
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()
	handlers.Init()
	handlers.JobCh = make(chan handlers.DeleteURLJob, 100)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			for job := range handlers.JobCh {
				handlers.DelFunc(job)
			}
		}()
	}
	if config.GetEnableHTTPS() {
		if err := srv.ListenAndServeTLS("cmd/shortener/certificate/localhost.crt", "cmd/shortener/certificate/localhost.key"); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}
	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")
}
