package slave_side

import (
	"fmt"
	//"vcbb/master_side"

	"vcbb/peer_list"
)

type Scheduler struct {
	maxJobCount     uint64
	runningJobList  []*Job
	peerList        *peer_list.PeerList
	seekReq         chan peer_list.MessageInfo
	TerminateSignal chan struct{}
}

func NewScheduler(maxjobcnt uint64, peerlist *peer_list.PeerList) *Scheduler {
	return &Scheduler{
		maxJobCount: maxjobcnt,
		peerList:    peerlist,
	}
}

func (this *Scheduler) Init() {
	this.runningJobList = make([]*Job, 0)
	this.seekReq = make(chan peer_list.MessageInfo, 10)
	this.peerList.AddChannel(peer_list.SeekParticipantReq, this.seekReq)
	this.TerminateSignal = make(chan struct{}, 1)
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
			//var reqobj master_side.ComputationReq

		case <-this.TerminateSignal:
			return
		}
	}
}

func (this *Scheduler) Stop() {

}
