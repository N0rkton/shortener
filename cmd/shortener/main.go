package main

import (
	"bufio"
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
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type body struct {
	URL string `json:"url"`
}
type response struct {
	Result string `json:"result"`
}

func isGzip(r *http.Request) io.Reader {
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, _ := gzip.NewReader(r.Body)
		defer gz.Close()
		return gz
	}
	return r.Body
}
func jsonIndexPage(w http.ResponseWriter, r *http.Request) {
	read := isGzip(r)
	var body body
	err := json.NewDecoder(read).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(body.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := generateRandomString()
	ok := db.AddURL(code, body.URL)
	if ok != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	writeToFile(code, body.URL)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	baseURL := os.Getenv("BASE_URL")
	var res response
	if baseURL == "" {
		res.Result = *b + "/" + code
	} else {
		res.Result = baseURL + "/" + code
	}
	resp, _ := json.Marshal(res)
	w.Write(resp)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	read := isGzip(r)
	s, err := io.ReadAll(read)
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
		w.Write([]byte(*b + "/" + code))
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
const addr = "http://localhost:8080"

var (
	a *string
	b *string
	f *string
)

func init() {
	a = flag.String("a", "localhost:8080", "Domain name")
	b = flag.String("b", addr, "port number")
	f = flag.String("f", "", "file path")
}
func readFromFile() {
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	var text map[string]string
	if fileStoragePath == "" {
		fileStoragePath = *f
	}
	if fileStoragePath != "" {
		file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			return
		}
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
	if fileStoragePath == "" {
		fileStoragePath = *f
	}
	if fileStoragePath != "" {
		file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return
		}
		defer file.Close()
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

var db storage.Storage

func main() {
	flag.Parse()
	db = storage.NewMemoryStorage()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", jsonIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/{id}", redirectTo).Methods(http.MethodGet)
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		log.Fatal(http.ListenAndServe(*a, gzipHandle(router)))
	} else {
		log.Fatal(http.ListenAndServe(serverAddress, gzipHandle(router)))
	}
}
