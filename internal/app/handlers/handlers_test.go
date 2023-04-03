package handlers

import (
	"bytes"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_indexPage(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		want    want
		request string
		body    string
	}{
		{name: "Positive",
			want:    want{code: 201},
			request: "http://localhost:8080/",
			body:    "http://xnewqaajckkrj9.biz/dtncu35",
		},
		{name: "Negative",
			want:    want{code: 400},
			request: "http://localhost:8080/",
			body:    "xnewqaajckkrj9.biz/dtncu35",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
			if err != nil {
				log.Fatal(err)
			}
			Init()
			fileStorage, _ = storage.NewFileStorage(*config.FileStoragePath)
			localMem = storage.NewStorageMock()
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(IndexPage)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func Test_jsonIndexPage(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name    string
		want    want
		request string
		body    []byte
	}{

		{name: "Positive",
			want:    want{code: 201},
			request: "http://localhost:8080/",
			body:    []byte(`{"url":"http://localhost:8080/BpLnf"}`),
		},
		{name: "Negative",
			want:    want{code: 400},
			request: "http://localhost:8080/",
			body:    []byte(`{"result":"http://localhost:8080/BpLnf"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
			var err error
			secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
			if err != nil {
				log.Fatal(err)
			}
			localMem = storage.NewStorageMock()
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(JSONIndexPage)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func Test_redirectTo(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		want    want
		request string
		code    string
		link    string
	}{

		{name: "Positive",
			want:    want{code: 307},
			request: "http://localhost:8080/ShortedURL",
		},
		{name: "Negative",
			want:    want{code: 400},
			request: "http://localhost:8080/aaaaaaaaaaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()
			localMem = storage.NewStorageMock()
			r := mux.NewRouter()
			r.HandleFunc("/{id}", RedirectTo)
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, tt.request, nil))
			result := w2.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func BenchmarkIndexPage(b *testing.B) {
	var err error
	secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	Init()
	for i := 0; i < b.N; i++ {
		localMem = storage.NewStorageMock()
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("http://xnewqaajckkrj9.biz/dtncu35"))
		request.Header.Set("Content-Type", "text/plain; charset=utf-8")
		w := httptest.NewRecorder()
		h := http.HandlerFunc(IndexPage)
		h(w, request)
		result := w.Result()
		assert.Equal(b, 201, result.StatusCode)
		defer result.Body.Close()
	}
}
func ExamplePingDB() {

}
func ExampleIndexPage() {
	var err error
	secret, err = hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}
	Init()

	localMem = storage.NewStorageMock()
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("http://xnewqaajckkrj9.biz/dtncu35"))
	request.Header.Set("Content-Type", "text/plain; charset=utf-8")
	w := httptest.NewRecorder()
	h := http.HandlerFunc(IndexPage)
	h(w, request)
	result := w.Result()
	//equal := assert.Equal(201, result.StatusCode)
	defer result.Body.Close()
}
