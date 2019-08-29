package master_side

import (
	"math/big"
	"vcbb/blockchain"
	"vcbb/log"
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
	CalculationContract *blockchain.CalculationContract
	AnswerDistribute    map[string][]types.Address
	ParticipantState    map[string]bool
	ContractStateUpdate chan *blockchain.Answer
	AnswerCnt           uint8
	MinAnswerCnt        uint8
	MaxAnswerCnt        uint8
	MaxAnswer           [][]string
	MaxAnswerHash       string
	logger              *log.LoggerInstance
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

func NewJob(id string, sch *Scheduler, schnode *scheduleNode, minAnsCnt uint8, fatherLogger *log.LoggerInstance /*dependencies []*Dependency, partitions []string, code, basetest, hardwarereq string*/) *Job {
	ret := &Job{
		ID:           id,
		State:        Preparing,
		Sch:          sch,
		PeerList:     sch.peerList.GetInstance(id),
		SchNode:      schnode,
		MinAnswerCnt: minAnsCnt,
		logger:       fatherLogger.GetSubInstance(log.Topic("Job " + id)),
		//Dependencies:       dependencies,
		//Partitions:         partitions,
		//Code:               code,
		//BaseTest:           basetest,
		//HarWareRequirement: hardwarereq,
	}
	ret.logger.Log("Try To Resolve Dependencies")
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
	ret.logger.Log("Job Construct OK")
	return ret
}

func (this *Job) Init() {
	this.ContractStateUpdate = make(chan *blockchain.Answer, 5)
	binfo := &blockchain.ContractDeployInfo{
		Value:    big.NewInt(130),
		GasLimit: uint64(4712388),
	}
	cinfo := &blockchain.CalculationContractDeployInfo{
		Id:               this.ID,
		St:               big.NewInt(0),
		Fund:             big.NewInt(100),
		ParticipantCount: uint8(2),
		Distribute:       [8]*big.Int{big.NewInt(20), big.NewInt(10), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
	}
	this.CalculationContract = blockchain.NewCalculationContract(this.Sch.bcHandler, this.ContractStateUpdate, binfo, cinfo, this.logger)
	this.AnswerDistribute = make(map[string][]types.Address)
}
