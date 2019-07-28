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
	ParticipantState    map[types.Address]PtState
	PartitionDistribute []uint64
	AnswerDistribute    map[string][]types.Address
	MetaDataReq         chan peer_list.MessageInfo
	ContractStateUpdate chan *blockchain.ComputationContractUpdate
	TerminateSignal     chan struct{}
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

func NewJob(id string, sch *Scheduler, schnode *scheduleNode /*dependencies []*Dependency, partitions []string, code, basetest, hardwarereq string*/) *Job {
	return &Job{
		ID:       id,
		State:    Preparing,
		Sch:      sch,
		PeerList: sch.peerList.GetInstance(id),
		SchNode:  schnode,
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
	this.ParticipantState = make(map[types.Address]PtState)
	this.PartitionDistribute = make([]uint64, len(this.SchNode.partitions))
	this.AnswerDistribute = make(map[string][]types.Address)
	this.MetaDataReq = make(chan peer_list.MessageInfo, 1)
	this.TerminateSignal = make(chan struct{}, 2)
	this.PeerList.AddChannel(peer_list.MetaDataReq, this.MetaDataReq)
	//this.PeerList.Init(nil, nil, this.MetaDataReq)
}
