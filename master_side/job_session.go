package master_side

import (
	"encoding/json"
	"fmt"
	"vcbb/blockchain"
	"vcbb/slave_side"
	"vcbb/types"
)

func (this *Job) StartSession(sch *Scheduler) error {
	this.ComputationContract = blockchain.NewComputationContract(sch.bcHandler)
	addr, err := this.ComputationContract.Start()
	if err != nil {
		return err
	}
	req, err := this.getComputationReq(addr)
	err = sch.peerList.BroadCast(req)
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
		req := <-this.MetaDataReq
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
		sch.peerList.SendMsgTo(req.From, meta)
	}
}

func (this *Job) handleContractStateUpdate(sch *Scheduler) {
	for {
		upd := <-this.ContractStateUpdate
		sch.peerList.UpdatePunishedPeers(upd.Punished)
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
	this.ComputationContract.Terminate()
	this.State = Finished
}
