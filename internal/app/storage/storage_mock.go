package storage

import (
	"errors"
)

type StorageMock struct {
}

func (m *StorageMock) AddURL(id string, code string, url string) error {
	return nil
}

func (m *StorageMock) GetURL(url string) (string, error) {

	if url != "ShortedURL" {
		return "", errors.New("not found")
	}
	return "SomeLongURL", nil
}
func (m *StorageMock) GetURLByID(id string) (map[string]string, error) {
	if id != "1" {
		return nil, errors.New("not found")
	}
	return nil, nil
}

func NewStorageMock() Storage {
	return &StorageMock{}
}
