package data

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
)

type Tracker struct {
	ID              string
	session         *DataTransportSession
	peerList        *peer_list.PeerListInstance
	infoReq         chan peer_list.MessageInfo
	TerminateSignal chan struct{}
}

func NewTracker(id string, session *DataTransportSession, pl *peer_list.PeerList) *Tracker {
	return &Tracker{
		ID:       id,
		session:  session,
		peerList: pl.GetInstance(id),
	}
}

func (this *Tracker) StartTracker() {
	this.infoReq = make(chan peer_list.MessageInfo, 1)
	this.peerList.AddChannel(peer_list.InfoReq, this.infoReq)
	this.TerminateSignal = make(chan struct{}, 1)
	go this.serve()
}

func (this *Tracker) serve() {
	for {
		select {
		case req := <-this.infoReq:
			state, err := this.session.GetState()
			if err != nil {
				continue
			}
			res := dataInfoRes{DataReceivers: state}
			resb, err := json.Marshal(res)
			if err != nil {
				continue
			}
			this.peerList.RemoteProcedureCall(req.From, this.ID, peer_list.InfoRes, resb)
			//this.peerList.SendMsgTo(req.From, resb)
		case <-this.TerminateSignal:
			return
		}
	}
}

func (this *Tracker) Terminate() {
	this.TerminateSignal <- *new(struct{})
}
