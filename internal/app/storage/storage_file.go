package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type FileStorage struct {
	memStorage Storage
	f          *os.File
}

func NewFileStorage(path string) (Storage, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()
	memDB := NewMemoryStorage()
	scanner := bufio.NewScanner(file)
	var text map[string]map[string]string
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &text); err != nil {
			return nil, fmt.Errorf("unable to unmarshal metric from file: %w", err)
		}
		for key, value := range text {
			for k, v := range value {
				memDB.AddURL(key, k, v)
			}
		}
	}
	return &FileStorage{
		memStorage: memDB,
		f:          file,
	}, nil
}

func (fs *FileStorage) AddURL(id string, code string, url string) error {
	fs.memStorage.AddURL(id, code, url)
	text, err := json.Marshal(map[string]map[string]string{id: {code: url}})
	if err != nil {
		return errors.New("json error")
	}
	text = append(text, '\n')
	fs.f.Write(text)
	return nil
}

func (fs *FileStorage) GetURL(url string) (string, error) {
	return fs.memStorage.GetURL(url)
}
func (fs *FileStorage) GetURLById(id string) (map[string]string, error) {
	return fs.memStorage.GetURLById(id)
}
