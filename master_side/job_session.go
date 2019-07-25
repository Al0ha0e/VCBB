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
	go this.HandleMetaDataReq(sch)
	go this.HandleContractStateUpdate(sch)
	return nil
}

//pack contract address,hardware requirement and basetest into a request message
func (this *Job) getComputationReq(addr types.Address) ([]byte, error) {
	ret := computationReq{ContractAddr: addr, Hardware: this.HarWareRequirement, BaseTest: this.BaseTest}
	retb, err := json.Marshal(ret)
	return retb, err
}

func (this *Job) getJobByAmount(amount uint64) ([]byte, error) {
	if amount == 0 {
		return nil, fmt.Errorf("amount of metadata must be positive")
	}
	if amount > uint64(len(this.Partitions)) {
		amount = uint64(len(this.Partitions))
	}
	ret := MetaDataRes{Code: this.Code}
	var dep []*JobMeta
	for _, meta := range this.Dependencies {
		dep = append(dep, meta.DependencyJobMeta)
	}
	ret.DependencyMeta = dep
	retb, err := json.Marshal(ret)
	return retb, err
}

func (this *Job) updateAnswer(map[string][]types.Address) (bool, error) {
	return true, nil
}

//reply to the peer who ask for the metadata to start a computation
//code and metadata
func (this *Job) HandleMetaDataReq(sch *Scheduler) {
	for true {
		req := <-this.MetaDataReq
		var reqobj slave_side.MetaDataReq
		err := json.Unmarshal(req.Content, &reqobj)
		if err != nil {
			continue
		}
		meta, err := this.getJobByAmount(reqobj.Amount)
		if err != nil {
			continue
		}
		sch.peerList.SendMsgTo(req.From, meta)
	}
}

func (this *Job) HandleContractStateUpdate(sch *Scheduler) {
	for true {
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
