package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type FileStorage struct {
	memStorage Storage
	f          *os.File
}

func NewFileStorage(path string) (Storage, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)
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
	var err error
	fs.f, err = os.OpenFile(fs.f.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", fs.f.Name(), err)
	}
	log.Println(fs.f.Name())
	fs.memStorage.AddURL(id, code, url)
	text, err := json.Marshal([]string{id, code, url})
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
func (fs *FileStorage) GetURLByID(id string) (map[string]string, error) {
	return fs.memStorage.GetURLByID(id)
}
func (fs *FileStorage) Del(id string, code string) {
	fs.memStorage.Del(id, code)
}
