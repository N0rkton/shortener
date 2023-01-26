package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			r.Body = gzipDecode(r)
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		r.Body = gzipDecode(r)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type body struct {
	URL string `json:"url"`
}
type response struct {
	Result string `json:"result"`
}

func gzipDecode(r *http.Request) io.ReadCloser {
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, _ := gzip.NewReader(r.Body)
		defer gz.Close()
		return gz
	}
	return r.Body
}
func jsonIndexPage(w http.ResponseWriter, r *http.Request) {
	var body body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(body.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	code := generateRandomString()
	ok := db.AddURL(code, body.URL)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	ok = fileStorage.AddURL(code, body.URL)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var res response
	res.Result = *config.baseURL + "/" + code
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("jsonIndexPage: encoding response:", err) // лучше это залоггировать тк это на нашей стороне проблема
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	s, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(string(s)) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	code := generateRandomString()
	ok := db.AddURL(code, string(s))
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	ok = fileStorage.AddURL(code, string(s))
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(*config.baseURL + "/" + code))
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["id"]

	link, ok := fileStorage.GetURL(shortLink)
	if link != "" {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	link, ok = db.GetURL(shortLink)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const urlLen = 5
const defaultBaseURL = "http://localhost:8080"

var config struct {
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
}

func init() {
	config.serverAddress = flag.String("a", "localhost:8080", "server address")
	config.baseURL = flag.String("b", defaultBaseURL, "base URL")
	config.fileStoragePath = flag.String("f", "", "file path")
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

var db storage.Storage
var fileStorage storage.Storage

func main() {
	flag.Parse()
	serverAddressEnv := os.Getenv("SERVER_ADDRESS")
	if serverAddressEnv != "" {
		config.serverAddress = &serverAddressEnv
	}
	baseURLEnv := os.Getenv("BASE_URL")
	if baseURLEnv != "" {
		config.baseURL = &baseURLEnv
	}
	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		config.fileStoragePath = &fileStoragePathEnv
	}
	fileStorage = storage.NewFileStorage(*config.fileStoragePath)
	db = storage.NewMemoryStorage()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", jsonIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/{id}", redirectTo).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(*config.serverAddress, gzipHandle(router)))
}
