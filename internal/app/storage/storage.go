// Package storage provides implementations for data storage functions.
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

// Storage an interface that defines the following methods:
type Storage interface {
	//AddURL - add new URL to storage, where id - user cookie, code - short URL, url - original URL.
	AddURL(id string, code string, url string) error
	//GetURL - returns original URL by shorted URL.
	GetURL(code string) (string, error)
	//GetURLByID - returns all shorted and original URLs by user.
	GetURLByID(id string) (map[string]string, error)
	//Del - deletes URL from storage.
	Del(id string, code string)
	//GetStats return amount of shorted urls and users
	GetStats() (urls int32, users int32, err error)
}
type storeInfo struct {
	cookie      string
	originalURL string
	deleted     bool
}

// MemoryStorage a struct that implements the Storage interface and stores data in the computer's memory.
type MemoryStorage struct {
	localMem map[string]storeInfo
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new MemoryStorage instance.
func NewMemoryStorage() Storage {
	return &MemoryStorage{localMem: make(map[string]storeInfo)}
}

// AddURL adds a new URL to the storage, where id is the user cookie, code is the short URL, and url is the original URL.
func (sm *MemoryStorage) AddURL(id string, code string, url string) error {
	sm.localMem[code] = storeInfo{cookie: id, originalURL: url, deleted: false}
	return nil
}

// GetURL returns the original URL by the shorted URL.
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

// GetURLByID returns all shorted and original URLs by user.
func (sm *MemoryStorage) GetURLByID(id string) (map[string]string, error) {
	resp := make(map[string]string)
	for k, v := range sm.localMem {
		if v.cookie == id {
			resp[k] = v.originalURL
		}
	}
	return resp, nil
}

// Del deletes URL from storage.
func (sm *MemoryStorage) Del(id string, code string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	link, ok := sm.localMem[code]
	if ok && link.cookie == id {
		sm.localMem[code] = storeInfo{cookie: id, originalURL: link.originalURL, deleted: true}
		return
	}
}

// GetStats - returns amount of shorted URLS and users
func (sm *MemoryStorage) GetStats() (urls int32, users int32, err error) {

	usersMap := make(map[string]int)
	for _, v := range sm.localMem {
		if !v.deleted {
			urls += 1
		}
		usersMap[v.cookie] = 1
	}
	users = int32(len(usersMap))
	return urls, users, nil
}
