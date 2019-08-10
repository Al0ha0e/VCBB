package slave_side

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/blockchain"
	"github.com/Al0ha0e/vcbb/msg"
	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
	"github.com/Al0ha0e/vcbb/vcfs"
)

type Scheduler struct {
	maxJobCount     uint64
	runningJobList  []*Job
	peerList        *peer_list.PeerList
	fileSystem      *vcfs.FileSystem
	bcHandler       blockchain.BlockChainHandler
	TerminateSignal chan struct{}
	jobError        chan jobRunTimeError
}

func NewScheduler(maxjobcnt uint64, peerlist *peer_list.PeerList, fs *vcfs.FileSystem, bchandler blockchain.BlockChainHandler) *Scheduler {
	return &Scheduler{
		maxJobCount: maxjobcnt,
		peerList:    peerlist,
		fileSystem:  fs,
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
	//TODO: CHECK CONTRACT
	//TODO: CHECK BASETEST&HARDWARE
	sess := NewJob(types.NewAddress("reqobj.Master"), "reqobj.ID", reqobj.BaseTest, reqobj.Hardware, this)
	go sess.StartSession(req)
}

func (this *Scheduler) Stop() {

}
