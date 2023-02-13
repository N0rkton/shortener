package main

import (
	"github.com/N0rkton/shortener/internal/app/handlers"
	"github.com/N0rkton/shortener/internal/config"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	handlers.Init()
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.IndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", handlers.JSONIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten/batch", handlers.Batch).Methods(http.MethodPost)
	router.HandleFunc("/ping", handlers.PingDB).Methods(http.MethodGet)
	router.HandleFunc("/{id}", handlers.RedirectTo).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", handlers.ListURL).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(config.GetServerAddress(), handlers.GzipHandle(router)))
}
