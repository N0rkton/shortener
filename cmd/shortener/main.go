package main

import (
	"encoding/json"
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

func jsonIndexPage(w http.ResponseWriter, r *http.Request) {
	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(b.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := generateRandomString()
	db.AddURL(code, b.URL)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var res response
	res.Result = addr + code
	resp, _ := json.Marshal(res)
	w.Write(resp)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
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
	ok := db.AddURL(code, string(s))
	if ok != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(addr + code))
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const urlLen = 5
const addr = "http://localhost:8080/"

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
	return err == nil && u.Host != ""
}

type body struct {
	URL string `json:"url"`
}
type response struct {
	Result string `json:"result"`
}

var db storage.Storage

func main() {
	db = storage.NewMemoryStorage()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", jsonIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/{id}", redirectTo).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
