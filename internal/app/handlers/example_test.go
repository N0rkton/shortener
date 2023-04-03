package handlers

import (
	"fmt"
	"github.com/N0rkton/shortener/internal/app/cookies"
	"io"
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

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	// Output:
	// 201
	// plain/text

}
func ExampleRedirectTo() {
	Init()
	localMem.AddURL("userID", "ABC123", "https://ya.ru")

	request := httptest.NewRequest(http.MethodGet, "http://example/ABC123", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(RedirectTo)
	h(w, request)
	resp := w.Result()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	// Output:
	// 307
	// plain/text
}
func ExampleListURL() {

	Init()
	val := generateRandomString(3)
	localMem.AddURL(val, "ABC123", "https://ya.ru")
	localMem.AddURL(val, "123ABC", "https://vk.com")
	request := httptest.NewRequest(http.MethodGet, "http://example/api/user/urls", nil)
	cookie := http.Cookie{
		Name:     "UserId",
		Value:    val,
		Path:     "/",
		HttpOnly: true,
		Secure:   false}

	w := httptest.NewRecorder()
	cookies.WriteEncrypted(w, cookie, secret)

	h := http.HandlerFunc(ListURL)
	h(w, request)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	// Output:
	// 204
	// text/plain; charset=utf-8
	// kk

}
