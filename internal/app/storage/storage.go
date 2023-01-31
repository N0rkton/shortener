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
	localMem map[string]map[string]string
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{localMem: make(map[string]map[string]string)}
}
func (sm *MemoryStorage) AddURL(id string, code string, url string) error {
	tmp := make(map[string]string)
	for k, v := range sm.localMem[id] {
		tmp[k] = v
	}
	tmp[code] = url
	sm.localMem[id] = tmp
	return nil
}

func (sm *MemoryStorage) GetURL(url string) (string, error) {
	for k := range sm.localMem {
		link, ok := sm.localMem[k][url]
		if ok {
			return link, nil
		}
	}
	return "", errors.New("not found")
}

func (sm *MemoryStorage) GetURLByID(id string) (map[string]string, error) {
	text, ok := sm.localMem[id]
	if ok {

		return text, nil
	}
	return nil, errors.New("not found")
}
