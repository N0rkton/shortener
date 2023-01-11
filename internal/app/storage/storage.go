package storage

import "errors"

type Storage interface {
	AddURL(code string, url string) error
	GetURL(url string) (string, error)
}
type Store struct {
	DB map[string]string
}

func (sm *Store) AddURL(code string, url string) error {
	sm.DB[code] = url
	return nil
}

func (sm *Store) GetURL(url string) (string, error) {
	link, ok := sm.DB[url]
	if !ok {
		return "", errors.New("not found")
	}
	return link, nil
}
