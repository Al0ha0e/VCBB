package master_side

import (
	"github.com/Al0ha0e/vcbb/blockchain"
	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

type JobState uint8

const (
	Preparing JobState = iota
	Running
	Finished
)

type PtState uint8

const (
	Unknown PtState = iota
	GotMeta
	Committed
)

type Job struct {
	ID                  string
	State               JobState
	Sch                 *Scheduler
	PeerList            *peer_list.PeerListInstance
	SchNode             *scheduleNode
	ComputationContract *blockchain.ComputationContract
	AnswerDistribute    map[string][]types.Address
	ContractStateUpdate chan *blockchain.ComputationContractUpdate
	AnswerCnt           uint8
	MinAnswerCnt        uint8
	MaxAnswerCnt        uint8
	MaxAnswer           string
	//Dependencies        []*Dependency//
	//Partitions          []string//
	//Code                string//
	//BaseTest            string//
	//HarWareRequirement  string//
}

type JobMeta struct {
	job              *Job
	node             *scheduleNode
	Contract         types.Address   `json:"contract"`
	Id               string          `json:"id"`
	Participants     []types.Address `json:"participants"`
	Partitions       []string        `json:"partitions"`
	PartitionAnswers []string        `json:"answers"`
	RootHash         string          `json:"root"`
}

func NewJob(id string, sch *Scheduler, schnode *scheduleNode, minAnsCnt uint8 /*dependencies []*Dependency, partitions []string, code, basetest, hardwarereq string*/) *Job {
	return &Job{
		ID:           id,
		State:        Preparing,
		Sch:          sch,
		PeerList:     sch.peerList.GetInstance(id),
		SchNode:      schnode,
		MinAnswerCnt: minAnsCnt,
		//Dependencies:       dependencies,
		//Partitions:         partitions,
		//Code:               code,
		//BaseTest:           basetest,
		//HarWareRequirement: hardwarereq,
	}
}

func (this *Job) Init() {
	this.ContractStateUpdate = make(chan *blockchain.ComputationContractUpdate, 1)
	this.ComputationContract = blockchain.NewComputationContract(this.Sch.bcHandler, this.ContractStateUpdate)
	this.AnswerDistribute = make(map[string][]types.Address)
}
