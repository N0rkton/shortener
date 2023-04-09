package cookies

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	value := "fgh"
	for i := 0; i < b.N; i++ {
		cookie := http.Cookie{
			Name:     "UserId",
			Value:    value,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		}
		request := httptest.NewRequest(http.MethodPost, "http://example/", strings.NewReader("http://xnewqaajckkrj9.biz/dtncu35"))
		request.AddCookie(&cookie)
		_, _ = Read(request, "UserID")
	}
}
func BenchmarkWrite(b *testing.B) {
	value := "fgh"
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		cookie := http.Cookie{
			Name:     "UserId",
			Value:    value,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		}
		_ = Write(w, cookie)
	}
}
