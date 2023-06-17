// Package storage provides functionality for storing URLs in text file.
package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// FileStorage - store data in txt file.
type FileStorage struct {
	memStorage Storage
	f          *os.File
}

// NewFileStorage function initializes a new FileStorage struct and reads the data from the provided file path.
// It creates a MemoryStorage struct to hold the data in memory, and adds each record to the MemoryStorage.
// It also marks the records that have been deleted, so that they can be skipped during future reads.
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
		if text[3] == "1" {
			memDB.Del(text[0], text[1])
		}
	}
	return &FileStorage{
		memStorage: memDB,
		f:          file,
	}, nil
}

// AddURL method adds a new record to the MemoryStorage and writes it to the file.
// It first opens the file in append mode, and then encodes the data as a JSON string and writes it to the file.
func (fs *FileStorage) AddURL(id string, code string, url string) error {
	var err error
	fs.f, err = os.OpenFile(fs.f.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", fs.f.Name(), err)
	}
	defer fs.f.Close()
	//log.Println(fs.f.Name())
	fs.memStorage.AddURL(id, code, url)
	text, err := json.Marshal([]string{id, code, url, "0"})
	if err != nil {
		return errors.New("json error")
	}
	text = append(text, '\n')
	fs.f.Write(text)
	return nil
}

// GetURL method looks up a URL by its code in the MemoryStorage.
func (fs *FileStorage) GetURL(code string) (string, error) {
	return fs.memStorage.GetURL(code)
}

// GetURLByID method looks up a URL by its ID in the MemoryStorage.
func (fs *FileStorage) GetURLByID(id string) (map[string]string, error) {
	return fs.memStorage.GetURLByID(id)
}

// Del method marks a record as deleted in the MemoryStorage and writes it to the file.
// It first checks whether the record exists in the MemoryStorage by calling the GetURL method.
// If the record has been deleted, it encodes the data as a JSON string and writes it to the file.
func (fs *FileStorage) Del(id string, code string) {
	fs.memStorage.Del(id, code)
	link, ok := fs.GetURL(code)
	if ok != nil {
		if errors.Is(ok, ErrDeleted) {
			var err error
			fs.f, err = os.OpenFile(fs.f.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				return
			}
			defer fs.f.Close()
			text, err := json.Marshal([]string{id, code, link, "1"})
			if err != nil {
				return
			}
			text = append(text, '\n')
			fs.f.Write(text)
		}
	}
}

// GetStats - returns amount of shorted URLS and users
func (fs *FileStorage) GetStats() (urls int32, users int32, err error) {
	return fs.memStorage.GetStats()
}
