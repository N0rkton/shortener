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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()
	var text []string
	memDB := NewMemoryStorage()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &text)
		if err != nil {
			return &FileStorage{
				memStorage: memDB,
				f:          file,
			}, fmt.Errorf("unable to unmarshall: %w", err)
		}
		memDB.AddURL(text[0], text[1], text[2])
	}
	return &FileStorage{
		memStorage: memDB,
		f:          file,
	}, nil
}

func (fs *FileStorage) AddURL(id string, code string, url string) error {
	fs.memStorage.AddURL(id, code, url)
	text, err := json.Marshal([]string{id, code, url})
	if err != nil {
		return errors.New("json error")
	}
	fs.f.Write(text)
	return nil
}

func (fs *FileStorage) GetURL(url string) (string, error) {
	return fs.memStorage.GetURL(url)
}
func (fs *FileStorage) GetURLByID(id string) (map[string]string, error) {
	return fs.memStorage.GetURLByID(id)
}
