package master_side

import (
	"encoding/json"
	"fmt"

	"vcbb/blockchain"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	fmt.Println("SESSION START", this.ID)
	this.Init()
	addr, err := this.CalculationContract.Start()
	if err != nil {
		fmt.Println("CONTRACT ERR", err)
		return err
	}
	fmt.Println("CONTRACT", addr)
	this.PeerList.AddCallBack("handleMetaDataReq", this.handleMetaDataReq)
	req, err := this.getComputationReq(addr)
	//err = this.PeerList.BroadCastRPC(this.ID, peer_list.SeekParticipantReq, req) //(req)
	_, err = sch.peerList.BroadCastRPC("handleSeekParticipantReq", req, 2)
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
		PartitionCnt: this.SchNode.partitionCnt,
		Hardware:     this.SchNode.hardwareRequirement,
		BaseTest:     this.SchNode.baseTest,
	}
	retb, err := json.Marshal(ret)
	return retb, err
}

func (this *Job) genPubKey(addr types.Address) string {
	return "SZH"
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
	res := msg.MetaDataRes{
		//TODO: SIGNATURE
		PublicKey:         this.genPubKey(req.From),
		Code:              this.SchNode.code,
		PartitionIdOffset: this.SchNode.partitionIDOffset,
		Inputs:            this.SchNode.input,
		DependencyMeta:    this.Dependencies,
	}
	resb, err := json.Marshal(res)
	this.PeerList.Reply(req, "handleMetaDataRes", resb)
}

func (this *Job) updateAnswer(newAnswer *blockchain.Answer) (bool, error) {
	fmt.Println("UPD ANS", newAnswer.Ans, newAnswer.AnsHash, newAnswer.Commiter)
	if this.ParticipantState[newAnswer.Commiter.ToString()] {
		return false, nil
	}
	k := newAnswer.AnsHash
	this.AnswerDistribute[k] = append(this.AnswerDistribute[k], newAnswer.Commiter)
	this.AnswerCnt++
	l := uint8(len(this.AnswerDistribute[k]))
	fmt.Println("L MAS_ANS_CNT", l, this.MaxAnswerCnt)
	if l > this.MaxAnswerCnt {
		this.MaxAnswerCnt = l
		this.MaxAnswer = newAnswer.Ans
		this.MaxAnswerHash = k
	}
	if this.AnswerCnt >= this.MinAnswerCnt && 2*this.MaxAnswerCnt > this.AnswerCnt {
		fmt.Println("CAN TERMINATE", this.ID)
		return true, nil
	}
	return false, nil
}

func (this *Job) handleContractStateUpdate(sch *Scheduler) {
	for {
		upd := <-this.ContractStateUpdate
		if len(upd.AnsHash) == 0 {
			this.PeerList.UpdatePunishedPeer(upd.Commiter)
			continue
		}
		canterminate, err := this.updateAnswer(upd)
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
	this.CalculationContract.Terminate() //TODO
	this.Sch.peerList.RemoveInstance(this.ID)
	ret := &JobMeta{
		node:             this.SchNode,
		Contract:         this.CalculationContract.Contract,
		Id:               this.ID,
		Participants:     this.AnswerDistribute[this.MaxAnswerHash],
		PartitionAnswers: this.MaxAnswer,
	}
	this.Sch.result <- ret
	fmt.Println("TERMINATED", this.ID)
}
