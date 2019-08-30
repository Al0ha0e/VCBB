package master_side

import (
	"encoding/json"
	"strconv"

	"vcbb/blockchain"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	this.logger.Log("Try To Start Job Session")
	this.Init()
	addr, err := this.CalculationContract.Start()
	if err != nil {
		this.logger.Err("Fail To Start Contract " + err.Error())
		return err
	}
	this.PeerList.AddCallBack("handleMetaDataReq", this.handleMetaDataReq)
	req, err := this.getComputationReq(addr)
	//err = this.PeerList.BroadCastRPC(this.ID, peer_list.SeekParticipantReq, req) //(req)
	_, err = sch.peerList.BroadCastRPC("handleSeekParticipantReq", req, 2)
	if err != nil {
		this.logger.Err("Broadcast Fail " + err.Error())
		return err
	}
	this.State = Running
	go this.handleContractStateUpdate(sch)
	this.logger.Log("Session Start")
	return nil
}

//pack contract address,hardware requirement and basetest into a request message
func (this *Job) getComputationReq(addr types.Address) ([]byte, error) {
	this.logger.Log("Try To Get Request")
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
	this.logger.Log("Receive MetaData Request From " + req.From.ToString() + " Session: " + req.FromSession)
	var reqobj msg.MetaDataReq
	//TODO: CHECK RESULT
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		this.logger.Err("Fail To Unmarshal MetaData Request " + err.Error())
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
	resb, _ := json.Marshal(res)
	this.PeerList.Reply(req, "handleMetaDataRes", resb)
}

func (this *Job) updateAnswer(newAnswer *blockchain.Answer) (bool, error) {
	this.logger.Log("Update Answer AnsHash: " + newAnswer.AnsHash + " From: " + newAnswer.Commiter.ToString())
	if this.ParticipantState[newAnswer.Commiter.ToString()] {
		return false, nil
	}
	k := newAnswer.AnsHash
	this.AnswerDistribute[k] = append(this.AnswerDistribute[k], newAnswer.Commiter)
	this.AnswerCnt++
	l := uint8(len(this.AnswerDistribute[k]))
	if l > this.MaxAnswerCnt {
		this.logger.Log("New Max Answer AnsHash: " + k + " Count: " + strconv.Itoa(int(this.MaxAnswerCnt)))
		this.MaxAnswerCnt = l
		this.MaxAnswer = newAnswer.Ans
		this.MaxAnswerHash = k
	}
	this.logger.Log("Update Answer OK AnswerHash: " + newAnswer.AnsHash + " AnswerCnt: " + strconv.Itoa(int(l)) + " TotalAnswerCnt: " + strconv.Itoa(int(this.AnswerCnt)))
	if this.AnswerCnt >= this.MinAnswerCnt && 2*this.MaxAnswerCnt > this.AnswerCnt {
		this.logger.Log("Terminate Condition OK")
		return true, nil
	}
	return false, nil
}

func (this *Job) handleContractStateUpdate(sch *Scheduler) {
	this.logger.Log("Start Watching COntract State At: " + this.CalculationContract.Contract.ToString())
	for {
		upd := <-this.ContractStateUpdate
		this.logger.Log("Receive Message From Contract Commiter: " + upd.Commiter.ToString())
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
	this.logger.Log("Try To Terminate SchNode ID: " + this.SchNode.id + " AnsHash: " + this.MaxAnswerHash + " AnswerCount: " + strconv.Itoa(int(this.MaxAnswerCnt)))
	this.State = Finished
	this.CalculationContract.Terminate()
	this.Sch.peerList.RemoveInstance(this.ID)
	ret := &JobMeta{
		node:             this.SchNode,
		Contract:         this.CalculationContract.Contract,
		Id:               this.ID,
		Participants:     this.AnswerDistribute[this.MaxAnswerHash],
		PartitionAnswers: this.MaxAnswer,
	}
	this.Sch.result <- ret
	this.logger.Log("Terminate OK")
}
