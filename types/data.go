package types

type data struct {
	V    []byte `json:"v"`
	Hash string `json:"hash"`
}

type SingleDataSource interface {
	GetHash() string
	GetValue() []byte
	GetSize() uint64
}

type DataSource interface {
	GetSingle(uint64) (SingleDataSource, error)
	GetHashList() []string
	GetHashSum(uint64, uint64) string
	GetTotalNum() uint64
}

type DataPack []SingleDataSource
