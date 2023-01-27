package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type FileStorage struct {
	fileStoragePath string
}

func NewFileStorage(path string) Storage {
	return &FileStorage{fileStoragePath: path}
}
func (f *FileStorage) AddURL(code string, url string) error {
	if f.fileStoragePath != "" {
		file, err := os.OpenFile(f.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return errors.New("cant open file")
		}
		defer file.Close()
		text, _ := json.Marshal(map[string]string{code: url})
		text = append(text, '\n')
		file.Write(text)
	}
	return nil
}

func (f *FileStorage) GetURL(url string) (string, error) {
	var text map[string]string
	db := NewMemoryStorage()
	if f.fileStoragePath != "" {
		file, err := os.OpenFile(f.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			return "", errors.New("cant open file")
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			json.Unmarshal(scanner.Bytes(), &text)
			for key, value := range text {
				db.AddURL(key, value)
			}
		}
		return db.GetURL(url)
	}
	return "", nil
}
