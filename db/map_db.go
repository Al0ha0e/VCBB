package db

import (
	"fmt"
	"sync"
)

type MapDB struct {
	kv     sync.Map
	lock   sync.Mutex
	rwLock sync.RWMutex
}

func (this *MapDB) GetMutex() *sync.Mutex {
	return &this.lock
}
func (this *MapDB) GetRWMutex() *sync.RWMutex {
	return &this.rwLock
}
func (this *MapDB) LockKey(key string) {
	this.lock.Lock()
}
func (this *MapDB) UnlockKey(key string) {
	this.lock.Unlock()
}
func (this *MapDB) RLockKey(key string) {
	this.rwLock.RLock()
}
func (this *MapDB) RUnlockKey(key string) {
	this.rwLock.RUnlock()
}
func (this *MapDB) WLockKey(key string) {
	this.rwLock.Lock()
}
func (this *MapDB) WUnlockKey(key string) {
	this.rwLock.Unlock()
}
func (this *MapDB) Get(key string) (interface{}, error) {
	value, ok := this.kv.Load(key)
	if !ok {
		return nil, fmt.Errorf("Key does not exist")
	}
	return value, nil
}
func (this *MapDB) Set(key string, value interface{}) error {
	this.kv.Store(key, value)
	return nil
}
func (this *MapDB) Has(key string) bool {
	_, ok := this.kv.Load(key)
	return ok
}
func (this *MapDB) Delete(key string) error {
	this.kv.Delete(key)
	return nil
}
func (this *MapDB) GetOrSet(key string, value interface{}) (interface{}, bool) {
	v, get := this.kv.LoadOrStore(key, value)
	return v, get
}
