package main

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/N0rkton/shortener/internal/app/cookies"
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
type idResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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
	var value string
	value, err := cookies.ReadEncrypted(r, "UserId", secret)
	if err != nil {
		value = generateRandomString(3)
		cookie := http.Cookie{
			Name:     "UserId",
			Value:    value,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		}
		err = cookies.WriteEncrypted(w, cookie, secret)
		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
	var body body
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(body.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	code := generateRandomString(urlLen)
	ok := localMem.AddURL(value, code, body.URL)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	if *config.fileStoragePath != "" {
		ok = fileStorage.AddURL(value, code, body.URL)
	}
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var res response
	res.Result = *config.baseURL + "/" + code
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("jsonIndexPage: encoding response:", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	var value string
	value, err := cookies.ReadEncrypted(r, "UserId", secret)
	if err != nil {
		value = generateRandomString(3)
		cookie := http.Cookie{
			Name:     "UserId",
			Value:    value,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		}
		err = cookies.WriteEncrypted(w, cookie, secret)
		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
	s, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidURL(string(s)) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	code := generateRandomString(urlLen)
	ok := localMem.AddURL(value, code, string(s))
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	if *config.fileStoragePath != "" {
		ok = fileStorage.AddURL(value, code, string(s))
	}
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
	var link string
	var ok error
	if *config.fileStoragePath != "" {
		link, ok = fileStorage.GetURL(shortLink)
	}
	if link != "" && ok == nil {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	link, ok = localMem.GetURL(shortLink)
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
func listURL(w http.ResponseWriter, r *http.Request) {
	var idR []idResponse
	var shortAndLongURL map[string]string
	var ok = errors.New("not found")
	value, err := cookies.ReadEncrypted(r, "UserId", secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	if fileStorage != nil {
		shortAndLongURL, ok = fileStorage.GetURLByID(value)
	}
	if ok == nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		for k, v := range shortAndLongURL {
			idR = append(idR, idResponse{ShortURL: *config.baseURL + "/" + k, OriginalURL: v})
		}
		if err := json.NewEncoder(w).Encode(idR); err != nil {
			log.Println("jsonIndexPage: encoding response:", err)
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
			return
		}
		return
	}
	shortAndLongURL, err = localMem.GetURLByID(value)
	if err != nil || shortAndLongURL == nil {
		http.Error(w, "no shorted urls", http.StatusNoContent)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	for k, v := range shortAndLongURL {
		idR = append(idR, idResponse{ShortURL: *config.baseURL + "/" + k, OriginalURL: v})
	}
	if err := json.NewEncoder(w).Encode(idR); err != nil {
		log.Println("jsonIndexPage: encoding response:", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}
func pingDB(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(context.Background(), *config.dbAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		//os.Exit(1)
	}
	defer conn.Close(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = conn.Ping(ctx); err != nil {
		http.Error(w, "unable to ping db", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

var secret []byte

const urlLen = 5
const defaultBaseURL = "http://localhost:8080"

var config struct {
	serverAddress   *string
	baseURL         *string
	fileStoragePath *string
	dbAddress       *string
}

func init() {
	config.serverAddress = flag.String("a", "localhost:8080", "server address")
	config.baseURL = flag.String("b", defaultBaseURL, "base URL")
	config.fileStoragePath = flag.String("f", "", "file path")
	config.dbAddress = flag.String("d", "", "data base connection address")
}

func generateRandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}

func isValidURL(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	return err == nil && u.Host != ""
}

var localMem storage.Storage
var fileStorage storage.Storage

func main() {
	flag.Parse()
	dbAddressEnv := os.Getenv("DATABASE_DSN")
	if dbAddressEnv != "" {
		config.dbAddress = &dbAddressEnv
	}
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
	var err error
	fileStorage, _ = storage.NewFileStorage(*config.fileStoragePath)
	localMem = storage.NewMemoryStorage()
	secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage).Methods(http.MethodPost)
	router.HandleFunc("/api/shorten", jsonIndexPage).Methods(http.MethodPost)
	router.HandleFunc("/{id}", redirectTo).Methods(http.MethodGet)
	router.HandleFunc("/api/user/urls", listURL).Methods(http.MethodGet)
	router.HandleFunc("/ping", pingDB).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(*config.serverAddress, gzipHandle(router)))
}
