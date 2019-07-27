package master_side

import (
	"encoding/json"
	"fmt"
	"vcbb/slave_side"
	"vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	this.Init()
	addr, err := this.ComputationContract.Start()
	if err != nil {
		return err
	}
	req, err := this.getComputationReq(addr)
	err = this.PeerList.BroadCast(req)
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
	ret := computationReq{ContractAddr: addr, Hardware: this.HarWareRequirement, BaseTest: this.BaseTest}
	retb, err := json.Marshal(ret)
	return retb, err
}

//get partitions by amount, [l,r) is the range where the partitions are most computed
func (this *Job) getJobByAmount(amount, l, r uint64) ([]byte, uint64, uint64, error) {
	tl, tr := l, r
	if amount == 0 {
		return nil, l, r, fmt.Errorf("amount of metadata must be positive")
	}
	pl := uint64(len(this.Partitions))
	if amount > pl {
		amount = pl
	}
	ret := MetaDataRes{Code: this.Code}
	for i := l - 1; i >= 0 && amount > 0; i++ {
		tl--
		amount--
		this.PartitionDistribute[i]++
		ret.Partitions = append(ret.Partitions, this.Partitions[i])
	}
	for i := r; i < pl && amount > 0; i++ {
		tr++
		amount--
		this.PartitionDistribute[i]++
		ret.Partitions = append(ret.Partitions, this.Partitions[i])
	}
	if amount > 0 {
		tl = l
		tr = l + 1
		for i := l; i < r && amount > 0; i++ {
			tr++
			amount--
			this.PartitionDistribute[i]++
			ret.Partitions = append(ret.Partitions, this.Partitions[i])
		}
	}
	var dep []*JobMeta
	for _, meta := range this.Dependencies {
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
			meta, l, r, err = this.getJobByAmount(reqobj.Amount, l, r)
			if err != nil {
				continue
			}
			this.ParticipantState[req.From] = GotMeta
			this.PeerList.SendMsgTo(req.From, meta)
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
		Partitions:   this.Partitions,
	}
	this.Sch.result <- ret
}
