package db

import "sync"

type DBIterator interface {
}

type DB interface {
	GetMutex() *sync.Mutex
	GetRWMutex() *sync.RWMutex
	LockKey(key string)
	UnlockKey(key string)
	RLockKey(key string)
	RUnlockKey(key string)
	WLockKey(key string)
	WUnlockKey(key string)
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Has(key string) bool
	Delete(key string) error
	GetOrSet(key string, value interface{}) (interface{}, bool)
}
