package master_side

import (
	"vcbb/blockchain"
	"vcbb/peer_list"
	"vcbb/types"
)

type JobState uint8

const (
	Preparing JobState = iota
	Running
	Finished
)

type Job struct {
	Id                  string
	State               JobState
	Dependencies        []*Dependency
	ComputationContract *blockchain.ComputationContract
	Partitions          []string
	PartitionDistribute []uint64
	AnswerDistribute    map[string][]types.Address
	MetaDataReq         chan peer_list.MessageInfo
	ContractStateUpdate chan *blockchain.ComputationContractUpdate
	Code                string
	BaseTest            string
	HarWareRequirement  string
}

type JobMeta struct {
	Id               string          `json:"id"`
	Participants     []types.Account `json:"participants"`
	Partitions       []string        `json:"partitions"`
	PartitionAnswers []string        `json:"answers"`
	RootHash         string          `json:"root"`
}
