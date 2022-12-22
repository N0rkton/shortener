package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func shorting() string {
	b := make([]byte, 5)
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

type Result struct {
	Link   string
	Code   string
	Status string
}

var db []Result

func indexPage(w http.ResponseWriter, r *http.Request) {
	result := Result{}
	if r.Method == "POST" {
		if !isValidURL(r.FormValue("s")) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			w.WriteHeader(400)
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.Code = shorting()
			w.WriteHeader(201)
			w.Header().Set("content-type", "text/plain")
			w.Write([]byte(result.Code))
			result.Status = "Сокращение было выполнено успешно"
		}
	}
	db = append(db, result)
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	b := 0
	for link := range db {
		if db[link].Code == vars["key"] {
			b = 1
			fmt.Fprintf(w, "<script>location='%s';</script>", db[link].Link)
			w.WriteHeader(307)
			w.Header().Set("Location", db[link].Link)
			break
		}
	}
	if b == 0 {
		w.WriteHeader(400)
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/to/{key}", redirectTo)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
