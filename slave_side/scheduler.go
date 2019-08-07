package slave_side

import (
	"encoding/json"
	"fmt"

	"github.com/Al0ha0e/vcbb/blockchain"
	"github.com/Al0ha0e/vcbb/master_side"
	"github.com/Al0ha0e/vcbb/peer_list"
)

type Scheduler struct {
	maxJobCount     uint64
	runningJobList  []*Job
	peerList        *peer_list.PeerList
	bcHandler       blockchain.BlockChainHandler
	seekReq         chan peer_list.MessageInfo
	TerminateSignal chan struct{}
	jobError        chan jobRunTimeError
}

func NewScheduler(maxjobcnt uint64, peerlist *peer_list.PeerList, bchandler blockchain.BlockChainHandler) *Scheduler {
	return &Scheduler{
		maxJobCount: maxjobcnt,
		peerList:    peerlist,
		bcHandler:   bchandler,
	}
}

func (this *Scheduler) Init() {
	this.runningJobList = make([]*Job, 0)
	this.seekReq = make(chan peer_list.MessageInfo, 10)
	this.peerList.AddChannel(peer_list.SeekParticipantReq, this.seekReq)
	this.TerminateSignal = make(chan struct{}, 1)
	this.jobError = make(chan jobRunTimeError, 10)
}

func (this *Scheduler) Run() {
	this.Init()
	go this.handleSeekParticipantReq()
}

func (this *Scheduler) handleSeekParticipantReq() {
	for {
		select {
		case req := <-this.seekReq:
			fmt.Println(req)
			var reqobj master_side.ComputationReq
			err := json.Unmarshal(req.Content, &reqobj)
			if err != nil {
				continue
			}
			//job := NewJob()
		case err := <-this.jobError:
			fmt.Println(err)
		case <-this.TerminateSignal:
			return
		}
	}
}

func (this *Scheduler) Stop() {

}
