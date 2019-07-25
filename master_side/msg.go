package master_side

import "vcbb/types"

type computationReq struct {
	ContractAddr types.Address `json:"address"`
	Hardware     string        `json:"hardware"`
	BaseTest     string        `json:"basetest"`
}

type MetaDataRes struct {
	Code           string     `json:"code"`
	Partitions     []string   `json:"partitions"`
	DependencyMeta []*JobMeta `json:"dependencyMeta"`
}
