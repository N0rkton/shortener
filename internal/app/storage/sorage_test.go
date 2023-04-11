package storage

import (
	"crypto/rand"
	"encoding/base32"
	"log"
	"os"
	"testing"
)

func generateRandomString(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}
func Benchmark_AddURL(b *testing.B) {
	file, err := os.OpenFile("test.txt", os.O_TRUNC, 0777)
	if err != nil {
		log.Print(err)
	}
	file.Close()
	fs, _ := NewFileStorage("test.txt")
	lm := NewMemoryStorage()

	b.Run("in memory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := generateRandomString(2)
			code := generateRandomString(4)
			url := generateRandomString(15)
			lm.AddURL(id, code, url)
		}
	})
	b.Run("in file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := generateRandomString(2)
			code := generateRandomString(4)
			url := generateRandomString(15)
			fs.AddURL(id, code, url)
		}
	})
}
func Benchmark_GetURL(b *testing.B) {
	file, err := os.OpenFile("test.txt", os.O_TRUNC, 0777)
	if err != nil {
		log.Print(err)
	}
	file.Close()
	fs, _ := NewFileStorage("test.txt")
	lm := NewMemoryStorage()
	fs.AddURL("1", "ABC", "http://example.com")
	lm.AddURL("1", "ABC", "http://example.com")

	b.Run("in memory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lm.GetURL("ABC")
		}
	})
	b.Run("in file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fs.GetURL("ABC")
		}
	})
}
func Benchmark_Del(b *testing.B) {
	file, err := os.OpenFile("test.txt", os.O_TRUNC, 0777)
	if err != nil {
		log.Print(err)
	}
	file.Close()
	fs, _ := NewFileStorage("test.txt")
	lm := NewMemoryStorage()

	b.Run("in memory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := generateRandomString(2)
			code := generateRandomString(4)
			url := generateRandomString(15)
			lm.AddURL(id, code, url)
			lm.Del(id, code)
		}
	})
	b.Run("in file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			id := generateRandomString(2)
			code := generateRandomString(4)
			url := generateRandomString(15)
			fs.AddURL(id, code, url)
			fs.Del(id, code)
		}
	})
}
