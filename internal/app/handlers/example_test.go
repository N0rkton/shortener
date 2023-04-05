package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
)

func ExampleIndexPage() {
	Init()
	request := httptest.NewRequest(http.MethodPost, "http://example/", strings.NewReader("http://xnewqaajckkrj9.biz/dtncu35"))
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w := httptest.NewRecorder()
	h := http.HandlerFunc(IndexPage)
	h(w, request)
	resp := w.Result()
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	// Output:
	// 201
	// plain/text

}
func ExampleRedirectTo() {
	Init()
	localMem.AddURL("", "ABC123", "http://ya.ru")
	r := mux.NewRouter()
	r.HandleFunc("/{id}", RedirectTo)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "http://example/ABC123", nil))
	resp := w.Result()
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)

	// Output:
	// 307

}
func ExampleListURL() {
	Init()
	r := mux.NewRouter()
	r.HandleFunc("/", IndexPage)
	r.HandleFunc("/api/user/urls", ListURL).Methods(http.MethodGet)
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8080", r))
	}()
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}
	rs, _ := c.Post("http://localhost:8080/", "text/plain; charset=utf-8", strings.NewReader("http://ya.ru"))
	defer rs.Body.Close()
	rs, _ = c.Get("http://localhost:8080/api/user/urls")
	defer rs.Body.Close()
	fmt.Println(rs.StatusCode)
	fmt.Println(rs.Header.Get("Content-Type"))

	// Output:
	// 200
	// application/json

}
