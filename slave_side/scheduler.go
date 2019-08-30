package slave_side

import (
	"encoding/json"

	"vcbb/blockchain"
	"vcbb/log"
	"vcbb/msg"
	"vcbb/peer_list"
	"vcbb/vcfs"
)

type Scheduler struct {
	peerList        *peer_list.PeerList
	fileSystem      *vcfs.FileSystem
	bcHandler       *blockchain.EthBlockChainHandler
	executer        Executer
	maxJobCount     uint64
	runningJobList  []*Job
	TerminateSignal chan struct{}
	jobError        chan jobRunTimeError
	logger          *log.LoggerInstance
}

func NewScheduler(
	maxjobcnt uint64,
	peerlist *peer_list.PeerList,
	fs *vcfs.FileSystem,
	bchandler *blockchain.EthBlockChainHandler,
	executer Executer,
	logSystem *log.LogSystem,
) *Scheduler {
	logSystem.Log("Create New Slave Scheduler")
	return &Scheduler{
		maxJobCount: maxjobcnt,
		peerList:    peerlist,
		fileSystem:  fs,
		executer:    executer,
		bcHandler:   bchandler,
		logger:      logSystem.GetInstance(log.Topic("Slave Scheduler")),
	}
}

func (this *Scheduler) Init() {
	this.runningJobList = make([]*Job, 0)
	this.TerminateSignal = make(chan struct{}, 1)
	this.jobError = make(chan jobRunTimeError, 10)
}

func (this *Scheduler) Run() {
	this.logger.Log("Start Scheduler")
	this.Init()
	this.peerList.AddCallBack("handleSeekParticipantReq", this.handleSeekParticipantReq)
}

func (this *Scheduler) handleSeekParticipantReq(req peer_list.MessageInfo) {
	this.logger.Log("Receive Job Request From: " + req.From.ToString() + " Session: " + req.FromSession)
	var reqobj msg.ComputationReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		this.logger.Err("Fail To Unmarshal Request")
		return
	}
	//fmt.Println("OBJ", reqobj)
	//TODO: CHECK CONTRACT
	//TODO: CHECK BASETEST&HARDWARE
	sess, err := NewJob(reqobj.ContractAddr, reqobj.Master, reqobj.Id, reqobj.BaseTest, reqobj.Hardware, this, reqobj.PartitionCnt, this.logger)
	if err != nil {
		return
	}
	go sess.StartSession(req)
}

func (this *Scheduler) Stop() {

}
