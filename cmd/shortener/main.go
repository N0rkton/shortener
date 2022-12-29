package main

import (
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const urlLen = 5

func generateRandomString() string {
	b := make([]byte, urlLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isValidURL(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

var db storage.Store

func indexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !isValidURL(string(s)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := generateRandomString()
		db.AddURL(code, string(s))
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("content-type", "plain/text")
		w.Write([]byte("http://localhost:8080/" + code))

	}
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["id"]
	link, ok := db.GetURL(shortLink)
	if ok != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	db = storage.NewStore()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/{id}", redirectTo)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
