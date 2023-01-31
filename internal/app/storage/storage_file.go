package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type FileStorage struct {
	memStorage Storage
	f          *os.File
}

func NewFileStorage(path string) (Storage, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()
	memDB := NewMemoryStorage()
	dat, _ := io.ReadAll(file)
	var text map[string]map[string]string
	if err := json.Unmarshal(dat, &text); err != nil {
		return nil, fmt.Errorf("unable to unmarshal metric from file: %w", err)
	}
	for key, value := range text {
		for k, v := range value {
			memDB.AddURL(key, k, v)
		}
	}

	return &FileStorage{
		memStorage: memDB,
		f:          file,
	}, nil
}

func (fs *FileStorage) AddURL(id string, code string, url string) error {
	fs.memStorage.AddURL(id, code, url)
	text, err := json.Marshal(fs.memStorage)
	if err != nil {
		return errors.New("json error")
	}
	fs.f.Seek(0, io.SeekStart)
	fs.f.Truncate(0)
	fs.f.Write(text)
	return nil
}

func (fs *FileStorage) GetURL(url string) (string, error) {
	return fs.memStorage.GetURL(url)
}
func (fs *FileStorage) GetURLByID(id string) (map[string]string, error) {
	return fs.memStorage.GetURLByID(id)
}
