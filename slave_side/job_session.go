package slave_side

import (
	"encoding/json"
	"fmt"

	"vcbb/blockchain"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

type jobRunTimeError struct {
	id  string
	err error
}

type Job struct {
	id                  string
	master              types.Address
	baseTest            string
	hardware            string
	code                string
	sch                 *Scheduler
	peerList            *peer_list.PeerListInstance
	partitionCnt        uint64
	calculationContract *blockchain.CalculationContract
}

func NewJob(contractAddress, master types.Address, id, baseTest, hardware string, sch *Scheduler, partitionCnt uint64) (*Job, error) {
	contract, err := blockchain.CalculationContractFromAddress(sch.bcHandler, contractAddress)
	if err != nil {
		return nil, err
	}
	return &Job{
		id:                  id,
		master:              master,
		baseTest:            baseTest,
		hardware:            hardware,
		sch:                 sch,
		partitionCnt:        partitionCnt,
		calculationContract: contract,
	}, nil
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
	this.peerList.Reply(req, "handleMetaDataReq", resb)
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
	//fmt.Println("OBJ", resobj)
	this.code = resobj.Code
	oksign := make(chan struct{}, 1)
	parts := make([]vcfs.FilePart, len(resobj.DependencyMeta))
	for i, meta := range resobj.DependencyMeta {
		parts[i].Peers = meta.Participants
		parts[i].Keys = meta.Keys
	}
	fmt.Println("PARTS", parts)
	go this.sch.fileSystem.FetchFiles(parts, oksign)
	<-oksign
	exeResultChan := make(chan *executeResult, 1)
	go this.sch.executer.Run(this.partitionCnt, resobj.PartitionIdOffset, resobj.Inputs, resobj.Code, exeResultChan)
	exeResult := <-exeResultChan
	fmt.Println(exeResult)
	// ERROR HANDLE
	this.calculationContract.Commit(nil, exeResult.result, "")
}

func (this *Job) terminate() {
	this.sch.peerList.RemoveInstance(this.id)
}
