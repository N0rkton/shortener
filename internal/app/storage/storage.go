package storage

import "errors"

type Storage interface {
	AddUrl(code string, url string)
	GetUrl(url string) (string, error)
	NewStore() Store
}
type Store struct {
	db map[string]string
}

var _ Store

func NewStore() Store {
	var ns Store
	ns.db = make(map[string]string)
	return ns
}
func (sm *Store) AddUrl(code string, url string) {
	sm.db[code] = url
}

func (sm *Store) GetURL(url string) (string, error) {
	link, ok := sm.db[url]
	if !ok {
		return "", errors.New("not found")
	}
	return link, nil
}
