package master_side

import "vcbb/types"

type ComputationReq struct {
	ContractAddr types.Address `json:"address"`
	PartitionCnt uint64        `json:"partitionCnt"`
	Hardware     string        `json:"hardware"`
	BaseTest     string        `json:"basetest"`
}

type MetaDataRes struct {
	PublicKey      string     `json:"publicKey"`
	Code           string     `json:"code"`
	Partitions     []string   `json:"partitions"`
	DependencyMeta []*JobMeta `json:"dependencyMeta"`
}
