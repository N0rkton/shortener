package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_indexPage(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name    string
		want    want
		request string
		body    string
	}{
		// TODO: Add test cases.
		{name: "Positive",
			want:    want{code: 201, response: "http://localhost:8080/" + shorting()},
			request: "http://localhost:8080/",
			body:    "http://xnewqaajckkrj9.biz/dtncu35",
		},
		{name: "Negative",
			want:    want{code: 400, response: "http://localhost:8080/" + shorting()},
			request: "http://localhost:8080/",
			body:    "xnewqaajckkrj9.biz/dtncu35",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(indexPage)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)

			//defer result.Body.Close()
			//resBody, err := io.ReadAll(result.Body)
			//if err != nil {
			//	t.Fatal(err)
			//}
			//if string(resBody) != tt.want.response {
			//	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			//	}

		})
	}
}

func Test_redirectTo(t *testing.T) {
	type want struct {
		code     int
		response string
		head     string
	}
	tests := []struct {
		name    string
		want    want
		request string
		body    string
	}{
		// TODO: Add test cases.
		{name: "Positive",
			want:    want{code: 400, response: "http://localhost:8080/" + shorting(), head: "Location"},
			request: "http://localhost:8080/",
			body:    "http://xnewqaajckkrj9.biz/dtncu35",
		},
		{name: "Negative",
			want:    want{code: 400, response: "http://localhost:8080/" + shorting()},
			request: "http://localhost:8080/aaaaaaaaaaa",
			body:    "http://xnewqaajckkrj9.biz/dtncu35",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain; charset=utf-8")
			w := httptest.NewRecorder()
			indexPage(w, req)
			request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/"+db[0].Code, nil)
			w2 := httptest.NewRecorder()

			fmt.Println(db)
			redirectTo(w2, request)
			result := w2.Result()
			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}
