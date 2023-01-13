package storage

import "errors"

type StorageMock struct {
}

func (m *StorageMock) AddURL(code string, url string) error {
	return nil
}

func (m *StorageMock) GetURL(url string) (string, error) {

	if url != "ShortedURL" {
		return "", errors.New("not found")
	}
	return "SomeLongURL", nil
}

func NewStorageMock() Storage {
	return &StorageMock{}
}
