package storage

import (
	"errors"
)

type Storage interface {
	AddURL(id string, code string, url string) error
	GetURL(url string) (string, error)
	GetURLByID(id string) (map[string]string, error)
}
type MemoryStorage struct {
	db map[string]map[string]string
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{db: make(map[string]map[string]string)}
}
func (sm *MemoryStorage) AddURL(id string, code string, url string) error {
	tmp := make(map[string]string)
	for k, v := range sm.db[id] {
		tmp[k] = v
	}
	tmp[code] = url
	sm.db[id] = tmp
	return nil
}

func (sm *MemoryStorage) GetURL(url string) (string, error) {
	for k := range sm.db {
		link, ok := sm.db[k][url]
		if ok {
			return link, nil
		}
	}
	return "", errors.New("not found")
}

func (sm *MemoryStorage) GetURLByID(id string) (map[string]string, error) {
	text, ok := sm.db[id]
	if ok {

		return text, nil
	}
	return nil, errors.New("not found")
}
