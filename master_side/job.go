package master_side

import (
	"vcbb/blockchain"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/types"
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
	Dependencies        []msg.JobMeta
	ComputationContract *blockchain.ComputationContract
	AnswerDistribute    map[string][]types.Address
	ContractStateUpdate chan *blockchain.ComputationContractUpdate
	AnswerCnt           uint8
	MinAnswerCnt        uint8
	MaxAnswerCnt        uint8
	MaxAnswer           [][]string
	MaxAnswerHash       string
	//Dependencies        []*Dependency//
	//Partitions          []string//
	//Code                string//
	//BaseTest            string//
	//HarWareRequirement  string//
}

type JobMeta struct {
	//job              *Job
	node             *scheduleNode
	Contract         types.Address   `json:"contract"`
	Id               string          `json:"id"`
	Participants     []types.Address `json:"participants"`
	PartitionAnswers [][]string      `json:"answers"`
	//Partitions       []string        `json:"partitions"`
	//RootHash         string          `json:"root"`
}

func NewJob(id string, sch *Scheduler, schnode *scheduleNode, minAnsCnt uint8 /*dependencies []*Dependency, partitions []string, code, basetest, hardwarereq string*/) *Job {
	ret := &Job{
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
	ret.Dependencies = make([]msg.JobMeta, 0)
	for _, v := range schnode.dependencies {
		keys := make([]string, len(v.keys))
		for i, id := range v.keys {
			pos := schnode.inputMap[id]
			keys[i] = schnode.input[pos.X][pos.Y]
		}
		dep := msg.JobMeta{
			Contract:     v.dependencyJobMeta.Contract,
			Participants: v.dependencyJobMeta.Participants,
			Keys:         keys,
		}
		ret.Dependencies = append(ret.Dependencies, dep)
	}
	return ret
}

func (this *Job) Init() {
	this.ContractStateUpdate = make(chan *blockchain.ComputationContractUpdate, 1)
	this.ComputationContract = blockchain.NewComputationContract(this.Sch.bcHandler, this.ContractStateUpdate)
	this.AnswerDistribute = make(map[string][]types.Address)
}
