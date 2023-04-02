// Package storage implements all data storage functions.
package storage

import (
	"errors"
	"sync"
)

// Module errors
var (
	ErrNotFound = errors.New("not found") // <- возвращаем когда урла совсем-совсем нет в базе
	ErrDeleted  = errors.New("deleted")   // <- возвращаем когда урл был, но удалили
)

type Storage interface {
	//AddURL - add new URL to storage, where id - user cookie, code - short URL, url - original URL.
	AddURL(id string, code string, url string) error
	//GetURL - returns original URL by shorted URL.
	GetURL(code string) (string, error)
	//GetURLByID - returns all shorted and original URLs by user.
	GetURLByID(id string) (map[string]string, error)
	//Del - deletes URL from storage.
	Del(id string, code string)
}
type storeInfo struct {
	cookie      string
	originalURL string
	deleted     bool
}

// MemoryStorage - store info in PC memory.
type MemoryStorage struct {
	localMem map[string]storeInfo
	mu       sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{localMem: make(map[string]storeInfo)}
}
func (sm *MemoryStorage) AddURL(id string, code string, url string) error {
	sm.localMem[code] = storeInfo{cookie: id, originalURL: url, deleted: false}
	return nil
}

func (sm *MemoryStorage) GetURL(code string) (string, error) {
	link, ok := sm.localMem[code]
	if !ok {
		return "", ErrNotFound
	}
	if link.deleted {
		return link.originalURL, ErrDeleted
	}
	return link.originalURL, nil
}

func (sm *MemoryStorage) GetURLByID(id string) (map[string]string, error) {
	resp := make(map[string]string)
	for k, v := range sm.localMem {
		if v.cookie == id {
			resp[k] = v.originalURL
		}
	}
	return resp, nil
}
func (sm *MemoryStorage) Del(id string, code string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	link, ok := sm.localMem[code]
	if ok && link.cookie == id {
		sm.localMem[code] = storeInfo{cookie: id, originalURL: link.originalURL, deleted: true}
		return
	}
}
