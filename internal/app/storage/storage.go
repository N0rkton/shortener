package storage

import "errors"

type Storage interface {
	AddURL(code string, url string) error
	GetURL(url string) (string, error)
}
type MemoryStorage struct {
	db map[string]string
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{db: make(map[string]string)}
}
func (sm *MemoryStorage) AddURL(code string, url string) error {
	sm.db[code] = url
	return nil
}

func (sm *MemoryStorage) GetURL(url string) (string, error) {
	link, ok := sm.db[url]
	if !ok {
		return "", errors.New("not found")
	}
	return link, nil
}
