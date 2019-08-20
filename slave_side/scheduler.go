package slave_side

import (
	"encoding/json"

	"vcbb/blockchain"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/vcfs"
)

type Scheduler struct {
	peerList        *peer_list.PeerList
	fileSystem      *vcfs.FileSystem
	bcHandler       *blockchain.EthBlockChainHandler
	executer        Executer
	maxJobCount     uint64
	runningJobList  []*Job
	TerminateSignal chan struct{}
	jobError        chan jobRunTimeError
}

func NewScheduler(maxjobcnt uint64, peerlist *peer_list.PeerList, fs *vcfs.FileSystem, bchandler *blockchain.EthBlockChainHandler, executer Executer) *Scheduler {
	return &Scheduler{
		maxJobCount: maxjobcnt,
		peerList:    peerlist,
		fileSystem:  fs,
		executer:    executer,
		bcHandler:   bchandler,
	}
}

func (this *Scheduler) Init() {
	this.runningJobList = make([]*Job, 0)
	this.TerminateSignal = make(chan struct{}, 1)
	this.jobError = make(chan jobRunTimeError, 10)
}

func (this *Scheduler) Run() {
	this.Init()
	this.peerList.AddCallBack("handleSeekParticipantReq", this.handleSeekParticipantReq)
}

func (this *Scheduler) handleSeekParticipantReq(req peer_list.MessageInfo) {
	var reqobj msg.ComputationReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	//fmt.Println("OBJ", reqobj)
	//TODO: CHECK CONTRACT
	//TODO: CHECK BASETEST&HARDWARE
	sess, err := NewJob(reqobj.ContractAddr, reqobj.Master, reqobj.Id, reqobj.BaseTest, reqobj.Hardware, this, reqobj.PartitionCnt)
	if err != nil {
		return
	}
	go sess.StartSession(req)
}

func (this *Scheduler) Stop() {

}
