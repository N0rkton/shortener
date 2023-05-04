// Package storage defining the behavior of a URL storage implementation.
package storage

import (
	"errors"
)

// StorageMock - struct to test Storage interface
type StorageMock struct {
}

// AddURL adds a new URL to the storage, where id is the user cookie, code is the short URL, and url is the original URL.
func (m *StorageMock) AddURL(id string, code string, url string) error {
	return nil
}

// GetURL returns the original URL by the shorted URL.
func (m *StorageMock) GetURL(url string) (string, error) {

	if url != "ShortedURL" {
		return "", ErrNotFound
	}
	return "SomeLongURL", nil
}

// GetURLByID returns all shorted and original URLs by user
func (m *StorageMock) GetURLByID(id string) (map[string]string, error) {
	if id != "1" {
		return nil, errors.New("not found")
	}
	return nil, nil
}

// Del delet URL
func (m *StorageMock) Del(id string, code string) {
}

// GetStats - returns amount of shorted URLS and users
func (m *StorageMock) GetStats() (urls int, users int, err error) {
	return 0, 0, nil
}

// NewStorageMock creates new mock instance.
func NewStorageMock() Storage {
	return &StorageMock{}
}
