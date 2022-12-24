package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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
		s, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if !isValidURL(string(s)) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			w.WriteHeader(400)
			result.Link = ""
		} else {
			result.Link = string(s)
			result.Code = shorting()
			result.Status = "Сокращение было выполнено успешно"
			db = append(db, result)
			w.WriteHeader(201)
			w.Header().Set("content-type", "plain/text")
			w.Write([]byte("http://localhost:8080/" + result.Code))

		}
	}

}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	b := 0
	for link := range db {
		if string(db[link].Code) == vars["id"] {
			b = 1

			w.WriteHeader(307)

			w.Header().Add("Location", db[link].Link)
			fmt.Print(w.Header().Values("Location"))
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
	router.HandleFunc("/{id}", redirectTo)
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
