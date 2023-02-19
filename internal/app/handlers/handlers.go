package handlers

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	conf "github.com/N0rkton/shortener/internal/app/config"
	"sync"

	"github.com/N0rkton/shortener/internal/app/cookies"
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var secret []byte
var localMem storage.Storage
var fileStorage storage.Storage
var db storage.Storage
var config conf.Cfg

func Init() {
	config = conf.NewConfig()
	var err error
	fileStorage, err = storage.NewFileStorage(*config.FileStoragePath)
	if err != nil {
		log.Println(err)
	}
	localMem = storage.NewMemoryStorage()
	db, err = storage.NewDBStorage(*config.DBAddress)
	if err != nil {
		log.Println(err)
	}
	secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func GzipHandle(next http.Handler) http.Handler {
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
type readBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
type respBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func gzipDecode(r *http.Request) io.ReadCloser {
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, _ := gzip.NewReader(r.Body)
		defer gz.Close()
		return gz
	}
	return r.Body
}
func JSONIndexPage(w http.ResponseWriter, r *http.Request) {
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
	if *config.FileStoragePath != "" {
		ok = fileStorage.AddURL(value, code, body.URL)
	}
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	if *config.DBAddress != "" {
		ok = db.AddURL(value, code, body.URL)
	}
	var pgErr *pgconn.PgError
	if errors.As(ok, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		link, ok2 := storage.GetShortURLByOrigin(*config.DBAddress, body.URL)
		if ok2 == nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusConflict)
			var res response
			res.Result = *config.BaseURL + "/" + link
			if err := json.NewEncoder(w).Encode(res); err != nil {
				log.Println("jsonIndexPage: encoding response:", err)
				http.Error(w, "unable to encode response", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var res response
	res.Result = *config.BaseURL + "/" + code
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println("jsonIndexPage: encoding response:", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}
func IndexPage(w http.ResponseWriter, r *http.Request) {
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
	if *config.FileStoragePath != "" {
		ok = fileStorage.AddURL(value, code, string(s))
	}
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	if *config.DBAddress != "" {
		ok = db.AddURL(value, code, string(s))
	}
	var pgErr *pgconn.PgError
	if errors.As(ok, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		link, ok2 := storage.GetShortURLByOrigin(*config.DBAddress, string(s))
		if link != "" && ok2 == nil {
			w.Header().Set("content-type", "plain/text")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(*config.BaseURL + "/" + link))
			return
		}
	}
	if ok != nil {
		http.Error(w, ok.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(*config.BaseURL + "/" + code))
}
func RedirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["id"]
	var link string
	var ok error
	if *config.DBAddress != "" {
		link, ok = db.GetURL(shortLink)
	}
	if ok != nil {
		if ok.Error() == "gone" {
			http.Error(w, ok.Error(), http.StatusGone)
			return
		}
	}
	if link != "" && ok == nil {
		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if *config.FileStoragePath != "" {
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
func ListURL(w http.ResponseWriter, r *http.Request) {
	var idR []idResponse
	var shortAndLongURL map[string]string
	var ok = errors.New("not found")
	value, err := cookies.ReadEncrypted(r, "UserId", secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	if *config.DBAddress != "" && fileStorage != nil {
		shortAndLongURL, ok = fileStorage.GetURLByID(value)
	}
	if ok == nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		for k, v := range shortAndLongURL {
			idR = append(idR, idResponse{ShortURL: *config.BaseURL + "/" + k, OriginalURL: v})
		}
		if err := json.NewEncoder(w).Encode(idR); err != nil {
			log.Println("jsonIndexPage: encoding response:", err)
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
			return
		}
		return
	}
	if *config.FileStoragePath != "" {
		shortAndLongURL, ok = fileStorage.GetURLByID(value)
	}
	if ok == nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		for k, v := range shortAndLongURL {
			idR = append(idR, idResponse{ShortURL: *config.BaseURL + "/" + k, OriginalURL: v})
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
		idR = append(idR, idResponse{ShortURL: *config.BaseURL + "/" + k, OriginalURL: v})
	}
	if err := json.NewEncoder(w).Encode(idR); err != nil {
		log.Println("jsonIndexPage: encoding response:", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}
func PingDB(w http.ResponseWriter, r *http.Request) {
	err := storage.Ping(*config.DBAddress)
	if err != nil {
		http.Error(w, "unable to ping db", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Batch(w http.ResponseWriter, r *http.Request) {
	var req []readBatch
	var resp []respBatch
	text, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(text, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for k := range req {
		if !isValidURL(req[k].OriginalURL) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		code := generateRandomString(urlLen)
		resp = append(resp, respBatch{req[k].CorrelationID, *config.BaseURL + "/" + code})
		ok := localMem.AddURL(req[k].CorrelationID, code, req[k].OriginalURL)
		if ok != nil {
			http.Error(w, ok.Error(), http.StatusBadRequest)
			return
		}
		if *config.FileStoragePath != "" {
			ok = fileStorage.AddURL(req[k].CorrelationID, code, req[k].OriginalURL)
		}
		if ok != nil {
			http.Error(w, ok.Error(), http.StatusBadRequest)
			return
		}
		if *config.DBAddress != "" {
			ok = db.AddURL(req[k].CorrelationID, code, req[k].OriginalURL)
		}
		if ok != nil {
			http.Error(w, ok.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("jsonIndexPage: encoding response:", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
		return
	}
}
func DeleteURL(w http.ResponseWriter, r *http.Request) {
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
	var text []string
	err = json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, "unable to decode body", http.StatusBadRequest)
	}
	log.Println(text)
	w.WriteHeader(http.StatusAccepted)
	wg := &sync.WaitGroup{}
	for _, v := range text {
		wg.Add(1)
		go func() {
			defer wg.Done()
			storage.Del(*config.DBAddress, v, value)
		}()
	}
	wg.Wait()
}

const urlLen = 5

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
