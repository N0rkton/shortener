package main

import (
	"bufio"
	"encoding/json"
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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
	ok := db.AddURL(code, b.URL)
	if ok != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	writeToFile(code, b.URL)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	baseURL := os.Getenv("BASE_URL")
	var res response
	if baseURL == "" {
		res.Result = addr + code

	} else {
		res.Result = baseURL + "/" + code
	}
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
	writeToFile(code, string(s))
	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		w.Write([]byte(addr + code))
		return
	}
	w.Write([]byte(baseURL + "/" + code))

}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["id"]
	readFromFile()
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

func readFromFile() {
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	var text map[string]string
	if fileStoragePath != "" {
		file, _ := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			json.Unmarshal(scanner.Bytes(), &text)
			for key, value := range text {
				db.AddURL(key, value)
			}
		}
	}
}
func writeToFile(code string, s string) {
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePath != "" {
		file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		defer file.Close()
		if err != nil {
			return
		}
		text, _ := json.Marshal(map[string]string{code: s})
		text = append(text, '\n')
		file.Write(text)
	}
}

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
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}
	db = storage.NewMemoryStorage()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", jsonIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/{id}", redirectTo).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
