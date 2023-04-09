package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "http://example/", strings.NewReader("http://xnewqaajckkrj9.biz/dtncu35"))
	h := http.HandlerFunc(IndexPage)
	h(w, request)
	resp := w.Result()
	defer resp.Body.Close()
	cookie := resp.Cookies()
	w2 := httptest.NewRecorder()
	request2 := httptest.NewRequest(http.MethodGet, "http://example/api/user/urls", nil)
	request2.AddCookie(cookie[0])
	http.SetCookie(w2, cookie[0])
	h2 := http.HandlerFunc(ListURL)
	h2(w2, request2)
	resp2 := w2.Result()
	fmt.Println(resp2.StatusCode)
	fmt.Println(resp2.Header.Get("Content-Type"))
	defer resp2.Body.Close()
	// Output:
	// 200
	// application/json

}
