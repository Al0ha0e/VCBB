package master_side

import (
	"encoding/json"
	"fmt"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/slave_side"
	"github.com/Al0ha0e/vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	this.Init()
	addr, err := this.ComputationContract.Start()
	if err != nil {
		return err
	}
	req, err := this.getComputationReq(addr)
	//err = this.PeerList.BroadCastRPC(this.ID, peer_list.SeekParticipantReq, req) //(req)
	err = sch.peerList.BroadCastRPC(peer_list.SeekParticipantReq, req)
	if err != nil {
		return err
	}
	this.State = Running
	go this.handleMetaDataReq(sch)
	go this.handleContractStateUpdate(sch)
	return nil
}

//pack contract address,hardware requirement and basetest into a request message
func (this *Job) getComputationReq(addr types.Address) ([]byte, error) {
	ret := ComputationReq{
		PartitionCnt: uint64(len(this.SchNode.partitions)),
		ContractAddr: addr,
		Hardware:     this.SchNode.hardwareRequirement,
		BaseTest:     this.SchNode.baseTest,
	}
	retb, err := json.Marshal(ret)
	return retb, err
}

func (this *Job) genPubKey(addr types.Address) string {
	return "SZH"
}

//get partitions by amount, [l,r) is the range where the partitions are most computed
func (this *Job) getJobByAmount(to types.Address, amount, l, r uint64) ([]byte, uint64, uint64, error) {
	tl, tr := l, r
	if amount == 0 {
		return nil, l, r, fmt.Errorf("amount of metadata must be positive")
	}
	pl := uint64(len(this.SchNode.partitions))
	if amount > pl {
		amount = pl
	}
	ret := MetaDataRes{PublicKey: this.genPubKey(to), Code: this.SchNode.code}
	for i := l - 1; i >= 0 && amount > 0; i++ {
		tl--
		amount--
		this.PartitionDistribute[i]++
		ret.Partitions = append(ret.Partitions, this.SchNode.partitions[i])
	}
	for i := r; i < pl && amount > 0; i++ {
		tr++
		amount--
		this.PartitionDistribute[i]++
		ret.Partitions = append(ret.Partitions, this.SchNode.partitions[i])
	}
	if amount > 0 {
		tl = l
		tr = l + 1
		for i := l; i < r && amount > 0; i++ {
			tr++
			amount--
			this.PartitionDistribute[i]++
			ret.Partitions = append(ret.Partitions, this.SchNode.partitions[i])
		}
	}
	var dep []*JobMeta
	for _, meta := range this.SchNode.dependencies {
		dep = append(dep, meta.DependencyJobMeta)
	}
	ret.DependencyMeta = dep
	retb, err := json.Marshal(ret)
	return retb, tl, tr, err
}

func (this *Job) updateAnswer(newAnswer map[string][]types.Address) (bool, error) {
	for k, v := range newAnswer {
		this.AnswerDistribute[k] = append(this.AnswerDistribute[k], v...)
	}
	return true, nil
}

//reply to the peer who ask for the metadata to start a computation
//code and metadata
func (this *Job) handleMetaDataReq(sch *Scheduler) {
	var l, r uint64
	for {
		select {
		case req := <-this.MetaDataReq:
			if this.ParticipantState[req.From] != Unknown {
				continue
			}
			var reqobj slave_side.MetaDataReq
			err := json.Unmarshal(req.Content, &reqobj)
			if err != nil {
				continue
			}
			var meta []byte
			meta, l, r, err = this.getJobByAmount(req.From, reqobj.Amount, l, r)
			if err != nil {
				continue
			}
			this.ParticipantState[req.From] = GotMeta
			this.PeerList.RemoteProcedureCall(req.From, this.ID, peer_list.MetaDataRes, meta) //SendMsgTo(req.From, meta)
		case <-this.TerminateSignal:
			return
		}

	}
}

func (this *Job) handleContractStateUpdate(sch *Scheduler) {
	for {
		select {
		case upd := <-this.ContractStateUpdate:
			this.PeerList.UpdatePunishedPeers(upd.Punished)
			canterminate, err := this.updateAnswer(upd.NewAnswer)
			if err != nil {
				continue
			}
			if canterminate {
				this.Terminate()
				break
			}
		case <-this.TerminateSignal: //BUG!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			return
		}

	}
}

func (this *Job) Terminate() {
	this.State = Finished
	this.ComputationContract.Terminate()
	this.PeerList.Close()
	this.TerminateSignal <- *new(struct{})
	this.TerminateSignal <- *new(struct{})
	close(this.TerminateSignal)
	var pt []types.Address
	for k, _ := range this.ParticipantState {
		pt = append(pt, k)
	}
	ret := &JobMeta{
		Id:           this.ID,
		Participants: pt,
		Partitions:   this.SchNode.partitions,
	}
	this.Sch.result <- ret
}
