package vcfs

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
