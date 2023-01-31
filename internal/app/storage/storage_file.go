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
	dat, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read: %w", err)
	}

	if err := json.Unmarshal(dat, &memDB); err != nil {
		return nil, fmt.Errorf("unable to unmarshal metric from file: %w", err)
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
	_, err = fs.f.Write(text)
	fmt.Println(string(text))
	if err != nil {
		print("json")
	}
	return nil
}

func (fs *FileStorage) GetURL(url string) (string, error) {
	return fs.memStorage.GetURL(url)
}
func (fs *FileStorage) GetURLByID(id string) (map[string]string, error) {
	return fs.memStorage.GetURLByID(id)
}
