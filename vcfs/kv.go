package vcfs

import "github.com/go-redis/redis"

type KVStore interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	CanSet(uint64) bool
}

type mapKVStore map[string][]byte

func neMapKVStore() *mapKVStore {
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
	client *redis.Client
}

func newRedisKVStore(addr string) *redisKVStore {
	return &redisKVStore{
		addr: addr,
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (this *redisKVStore) Get(key string) ([]byte, error) {
	v, err := this.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

func (this *redisKVStore) Set(key string, value []byte) error {
	err := this.client.Set(key, string(value), 0).Err()
	return err
}

func (this *redisKVStore) CanSet(uint64) bool {
	return true
}
