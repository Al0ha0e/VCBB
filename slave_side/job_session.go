package slave_side

import (
	"encoding/json"
	"math/big"

	"golang.org/x/crypto/sha3"

	"vcbb/blockchain"
	"vcbb/log"
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
	logger              *log.LoggerInstance
}

func NewJob(contractAddress, master types.Address, id, baseTest, hardware string, sch *Scheduler, partitionCnt uint64, fatherLogger *log.LoggerInstance) (*Job, error) {
	logger := fatherLogger.GetSubInstance(log.Topic(id))
	logger.Log("Try To Create New Job Session With Master " + master.ToString())
	contract, err := blockchain.CalculationContractFromAddress(sch.bcHandler, contractAddress, logger)
	if err != nil {
		logger.Err("Fail To Create Contrcat From Address " + contractAddress.ToString())
		return nil, err
	}
	logger.Log("Create Contract OK")
	return &Job{
		id:                  id,
		master:              master,
		baseTest:            baseTest,
		hardware:            hardware,
		sch:                 sch,
		partitionCnt:        partitionCnt,
		calculationContract: contract,
		logger:              logger,
	}, nil
}

func (this *Job) Init() {
	this.peerList = this.sch.peerList.GetInstance(this.id)
}

func (this *Job) StartSession(req peer_list.MessageInfo) {
	this.logger.Log("Start Session")
	this.Init()
	this.peerList.AddCallBack("handleMetaDataRes", this.handleMetaDataRes)
	res := msg.MetaDataReq{
		Result: "TODO",
	}
	resb, _ := json.Marshal(res)
	this.peerList.Reply(req, "handleMetaDataReq", resb)
}

func (this *Job) handleMetaDataRes(res peer_list.MessageInfo) {
	this.logger.Log("Receive Metadata From: " + res.From.ToString() + " Session: " + res.FromSession)
	if res.From != this.master {
		this.logger.Err("Wrong Metadata Origin")
		return
	}
	var resobj msg.MetaDataRes
	err := json.Unmarshal(res.Content, &resobj)
	if err != nil {
		this.logger.Err("Fail To Unmarshal Metadata " + err.Error())
		this.terminate()
		this.sch.jobError <- jobRunTimeError{
			id:  this.id,
			err: err,
		}
		return
	}
	this.logger.Log("Metadata Code " + resobj.Code)
	this.code = resobj.Code
	oksign := make(chan struct{}, 1)
	parts := make([]vcfs.FilePart, len(resobj.DependencyMeta))
	for i, meta := range resobj.DependencyMeta {
		parts[i].Peers = meta.Participants
		parts[i].Keys = meta.Keys
	}
	go this.sch.fileSystem.FetchFiles(parts, oksign)
	<-oksign
	this.logger.Log("File Fetch OK, Try To Execute")
	exeResultChan := make(chan *executeResult, 1)
	go this.sch.executer.Run(this.partitionCnt, resobj.PartitionIdOffset, resobj.Inputs, resobj.Code, exeResultChan)
	exeResult := <-exeResultChan
	this.logger.Log("Execute OK Try To Set FileInfo")
	// ERROR HANDLE
	var sum string
	for _, str := range exeResult.result {
		for _, str2 := range str {
			this.sch.fileSystem.SetInfo(str2)
			sum += str2
		}
	}
	sumb := make([]byte, 64)
	sha3.ShakeSum256(sumb, []byte(sum))
	this.logger.Log("Answer Hash " + string(sumb))
	info := &blockchain.ContractDeployInfo{
		Value:    big.NewInt(100),
		GasLimit: uint64(4712388),
	}
	this.logger.Log("Try To Commit Answer")
	this.calculationContract.Commit(info, exeResult.result, string(sumb))
	this.logger.Log("Contract Commit OK")
}

func (this *Job) terminate() {
	this.sch.peerList.RemoveInstance(this.id)
}
