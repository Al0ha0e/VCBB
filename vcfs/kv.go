package vcfs

type KVStore interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	CanSet(uint64) bool
}
