package vcfs

import (
	"sync"

	"github.com/gomodule/redigo/redis"
)

type KVStore interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	CanSet(uint64) bool
}

type mapKVStore map[string][]byte

func NewMapKVStore() *mapKVStore {
	var ret mapKVStore = make(map[string][]byte)
	return &ret
}

func (this *mapKVStore) Get(key string) ([]byte, error) {
	return (*this)[key], nil
}

func (this *mapKVStore) Set(key string, value []byte) error {
	(*this)[key] = value
	return nil
}

func (this *mapKVStore) CanSet(uint64) bool {
	return true
}

type redisKVStore struct {
	addr   string
	client redis.Conn
	lock   sync.Mutex
}

func NewRedisKVStore(addr string, db int) (*redisKVStore, error) {
	c, err := redis.Dial("tcp", addr, redis.DialDatabase(db))
	if err != nil {
		return nil, err
	}
	return &redisKVStore{
		addr:   addr,
		client: c,
		/*redis.NewClient(&redis.Options{
			Addr:     addr,
			DB:       db,
			PoolSize: 12,
		}),*/
	}, nil
}

func (this *redisKVStore) Get(key string) ([]byte, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	v, err := redis.Bytes(this.client.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (this *redisKVStore) Set(key string, value []byte) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, err := this.client.Do("SET", key, string(value))
	return err
}

func (this *redisKVStore) CanSet(uint64) bool {
	return true
}
