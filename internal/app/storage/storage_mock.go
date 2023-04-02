package storage

import (
	"errors"
)

// StorageMock - struct to test Storage interface
type StorageMock struct {
}

func (m *StorageMock) AddURL(id string, code string, url string) error {
	return nil
}

func (m *StorageMock) GetURL(url string) (string, error) {

	if url != "ShortedURL" {
		return "", ErrNotFound
	}
	return "SomeLongURL", nil
}
func (m *StorageMock) GetURLByID(id string) (map[string]string, error) {
	if id != "1" {
		return nil, errors.New("not found")
	}
	return nil, nil
}
func (m *StorageMock) Del(id string, code string) {
}

func NewStorageMock() Storage {
	return &StorageMock{}
}
