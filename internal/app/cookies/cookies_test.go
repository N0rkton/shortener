package cookies

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//	value := "fgh"
		//	cookie := http.Cookie{
		//		Name:     "UserId",
		//		Value:    value,
		//		Path:     "/",
		//		HttpOnly: true,
		//		Secure:   false,
		//	}
		//var r *http.Request

		//Read(r, "UserID")
	}
}
func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//		var w http.ResponseWriter
		//		value := "fgh"
		//		cookie := http.Cookie{
		//			Name:     "UserId",
		//			Value:    value,
		//			Path:     "/",
		//			HttpOnly: true,
		//			Secure:   false,
		//		}
		_ = httptest.NewRecorder()

	}
}
