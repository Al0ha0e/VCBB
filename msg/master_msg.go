package msg

import "github.com/Al0ha0e/vcbb/types"

type JobMeta struct {
	Contract         types.Address   `json:"contract"`
	Participants     []types.Address `json:"participants"`
	Partitions       []string        `json:"partitions"`
	PartitionAnswers []string        `json:"answers"`
}

type ComputationReq struct {
	Id           string        `json:"id"`
	Master       types.Address `json:"master"`
	ContractAddr types.Address `json:"address"`
	PartitionCnt uint64        `json:"partitionCnt"`
	Hardware     string        `json:"hardware"`
	BaseTest     string        `json:"basetest"`
}

type MetaDataRes struct {
	PublicKey      string     `json:"publicKey"`
	Signature      string     `json:"signature"`
	Code           string     `json:"code"`
	Partitions     []string   `json:"partitions"`
	DependencyMeta []*JobMeta `json:"dependencyMeta"`
}

type JobResult struct {
	ContractAddr types.Address `json:"address"`
}
