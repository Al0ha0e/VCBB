package vcfs

type KVStore interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
}
