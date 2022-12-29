package main

import (
	"github.com/N0rkton/shortener/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
			db = storage.NewStore()
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(indexPage)
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
			request: "http://localhost:8080/AAAAA",
			code:    "AAAAA",
			link:    "https://ya.ru",
		},
		{name: "Negative",
			want:    want{code: 400},
			request: "http://localhost:8080/aaaaaaaaaaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db = storage.NewStore()
			db.AddUrl(tt.code, tt.link)

			r := mux.NewRouter()
			r.HandleFunc("/{id}", redirectTo)
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, tt.request, nil))
			result := w2.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)
			defer result.Body.Close()
		})
	}
}
