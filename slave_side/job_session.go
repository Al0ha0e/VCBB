package slave_side

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/msg"
	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
	"github.com/Al0ha0e/vcbb/vcfs"
)

type jobRunTimeError struct {
	id  string
	err error
}

type Job struct {
	id       string
	master   types.Address
	baseTest string
	hardware string
	sch      *Scheduler
	peerList *peer_list.PeerListInstance
}

func NewJob(master types.Address, id, baseTest, hardware string, sch *Scheduler) *Job {
	return &Job{
		id:       id,
		master:   master,
		baseTest: baseTest,
		hardware: hardware,
		sch:      sch,
	}
}

func (this *Job) Init() {
	this.peerList = this.sch.peerList.GetInstance(this.id)
}

func (this *Job) StartSession(req peer_list.MessageInfo) {
	this.Init()
	this.peerList.AddCallBack("handleMetaDataRes", this.handleMetaDataRes)
	res := msg.MetaDataReq{
		Result: "TODO",
	}
	resb, _ := json.Marshal(res)
	this.peerList.Reply(req, "", resb)
}

func (this *Job) handleMetaDataRes(res peer_list.MessageInfo) {
	if res.From != this.master {
		return
	}
	var resobj msg.MetaDataRes
	err := json.Unmarshal(res.Content, &resobj)
	if err != nil {
		this.terminate()
		this.sch.jobError <- jobRunTimeError{
			id:  this.id,
			err: err,
		}
		return
	}
	oksign := make(chan struct{}, 1)
	parts := make([]vcfs.FilePart, len(resobj.DependencyMeta))
	for i, meta := range resobj.DependencyMeta {
		parts[i].peers = meta.Participants
		parts[i].keys = meta.PartitionAnswers
	}
	this.sch.fileSystem.FetchFiles(parts, oksign)
	<-oksign
}

func (this *Job) terminate() {
	this.sch.peerList.RemoveInstance(this.id)
}
