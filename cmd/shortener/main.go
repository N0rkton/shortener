package main

import (
	_ "net/http/pprof"

	"github.com/N0rkton/shortener/internal/app/config"
	"github.com/N0rkton/shortener/internal/app/handlers"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

const workerCount = 10

func main() {
	handlers.Init()
	handlers.JobCh = make(chan handlers.DeleteURLJob, 100)
	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range handlers.JobCh {
				handlers.DelFunc(job)
			}
		}()
	}
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", handlers.JSONIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten/batch", handlers.Batch).Methods(http.MethodPost)
	router.HandleFunc("/ping", handlers.PingDB).Methods(http.MethodGet)
	router.HandleFunc("/{id}", handlers.RedirectTo).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", handlers.ListURL).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", handlers.DeleteURL).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(config.GetServerAddress(), handlers.GzipHandle(router)))

}
