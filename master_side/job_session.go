package master_side

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/msg"
	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	this.Init()
	addr, err := this.ComputationContract.Start()
	if err != nil {
		return err
	}
	this.PeerList.AddCallBack("handleMetaDataReq", this.handleMetaDataReq)
	req, err := this.getComputationReq(addr)
	//err = this.PeerList.BroadCastRPC(this.ID, peer_list.SeekParticipantReq, req) //(req)
	_, err = sch.peerList.BroadCastRPC("", req, 2)
	if err != nil {
		return err
	}
	this.State = Running
	go this.handleContractStateUpdate(sch)
	return nil
}

//pack contract address,hardware requirement and basetest into a request message
func (this *Job) getComputationReq(addr types.Address) ([]byte, error) {
	ret := msg.ComputationReq{
		Id:           this.ID,
		Master:       this.Sch.peerList.Address,
		ContractAddr: addr,
		PartitionCnt: uint64(len(this.SchNode.partitions)),
		Hardware:     this.SchNode.hardwareRequirement,
		BaseTest:     this.SchNode.baseTest,
	}
	retb, err := json.Marshal(ret)
	return retb, err
}

func (this *Job) genPubKey(addr types.Address) string {
	return "SZH"
}

func (this *Job) updateAnswer(newAnswer map[string][]types.Address) (bool, error) {
	for k, v := range newAnswer {
		this.AnswerDistribute[k] = append(this.AnswerDistribute[k], v...)
		this.AnswerCnt += uint8(len(v))
		l := uint8(len(this.AnswerDistribute[k]))
		if l > this.MaxAnswerCnt {
			this.MaxAnswerCnt = l
			this.MaxAnswer = k
		}
	}
	if this.AnswerCnt >= this.MinAnswerCnt && 2*this.MaxAnswerCnt > this.AnswerCnt {
		return true, nil
	}
	return false, nil
}

//reply to the peer who ask for the metadata to start a computation
//code and metadata
func (this *Job) handleMetaDataReq(req peer_list.MessageInfo) {
	var reqobj msg.MetaDataReq
	//TODO: CHECK RESULT
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	deps := make([]*msg.JobMeta, len(this.SchNode.dependencies))
	for i, dep := range this.SchNode.dependencies {
		depmeta := dep.DependencyJobMeta
		deps[i] = &msg.JobMeta{
			Contract:         depmeta.Contract,
			Participants:     depmeta.Participants,
			Partitions:       depmeta.Partitions,
			PartitionAnswers: depmeta.PartitionAnswers,
		}
	}
	res := msg.MetaDataRes{
		//TODO: SIGNATURE
		PublicKey:      this.genPubKey(req.From),
		Code:           this.SchNode.code,
		Partitions:     this.SchNode.partitions,
		DependencyMeta: deps,
	}
	resb, err := json.Marshal(res)
	this.PeerList.Reply(req, "", resb)
}

func (this *Job) handleContractStateUpdate(sch *Scheduler) {
	for {
		upd := <-this.ContractStateUpdate
		this.PeerList.UpdatePunishedPeers(upd.Punished)
		canterminate, err := this.updateAnswer(upd.NewAnswer)
		if err != nil {
			continue
		}
		if canterminate {
			this.Terminate()
			break
		}

	}
}

func (this *Job) Terminate() {
	this.State = Finished
	this.ComputationContract.Terminate() //TODO
	this.Sch.peerList.RemoveInstance(this.ID)
	var pt []types.Address
	ret := &JobMeta{
		Id:           this.ID,
		Participants: pt,
		Partitions:   this.SchNode.partitions,
	}
	this.Sch.result <- ret
}
