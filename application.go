package main

import (
	"vcbb/blockchain"
	"vcbb/log"
	"vcbb/master_side"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/slave_side"
	"vcbb/types"
	"vcbb/vcfs"
)

type Application struct {
	Account           *types.Account
	NetService        net.NetService
	PeerList          *peer_list.PeerList
	FileSystem        *vcfs.FileSystem
	BlockchainHandler *blockchain.EthBlockChainHandler
	MasterClient      *master_side.MasterClient
	SlaveSide         *slave_side.Scheduler
	LogSystem         *log.LogSystem
}

func NewApplication(account, privateKey, inaddr, outaddr, dbaddr, bcurl string, db int) (ret *Application, err error) {
	ret = &Application{}
	logSystem, err := log.NewLogSystem("")
	if err != nil {
		return nil, err
	}
	ret.LogSystem = logSystem
	logSystem.Log("Try To Construct Application")
	acco := types.NewAddress(account)
	ret.Account = types.NewAccount(acco, privateKey)
	logSystem.Log("Try To Construct Net Service")
	ret.NetService, err = net.NewUDPNetService(inaddr, outaddr)
	if err != nil {
		logSystem.Err("Fail To Construct NetService " + err.Error())
		return nil, err
	}
	logSystem.Log("Try To Construct PeerList")
	ret.PeerList = peer_list.NewPeerList(acco, ret.NetService, logSystem)
	logSystem.Log("Try To Construct KVStore")
	engine, err := vcfs.NewRedisKVStore(dbaddr, db)
	if err != nil {
		logSystem.Err("Fail To Construct LogSystem " + err.Error())
		return nil, err
	}
	logSystem.Log("Try To Construct FileSystem")
	ret.FileSystem = vcfs.NewFileSystem(engine, ret.PeerList)
	logSystem.Log("Try To Construct Blockchain Handler")
	ret.BlockchainHandler, err = blockchain.NewEthBlockChainHandler(bcurl, ret.Account)
	if err != nil {
		logSystem.Err("Fail To Construct Blockchain Handler " + err.Error())
		return nil, err
	}
	return ret, err
}

func (this *Application) Run() {
	this.LogSystem.Log("Application Start")
	this.FileSystem.Serve()
	this.NetService.Run()
	this.PeerList.Run()
}

func (this *Application) StartMasterClient(client string) error {
	var err error
	this.LogSystem.Log("Try To Construct Master Client")
	this.MasterClient, err = master_side.NewMasterClient(this.PeerList, this.FileSystem, this.BlockchainHandler, this.LogSystem)
	if err != nil {
		this.LogSystem.Err("Fail To Construct Master Client " + err.Error())
		return err
	}
	this.MasterClient.Run(client)
	this.LogSystem.Log("Master Client Start")
	return nil
}

func (this *Application) StartSlaveSide(executer string) {
	this.LogSystem.Log("Try To Construct Slave Scheduler")
	exe := slave_side.NewPyExecuter(executer, this.LogSystem)
	this.SlaveSide = slave_side.NewScheduler(10, this.PeerList, this.FileSystem, this.BlockchainHandler, exe, this.LogSystem)
	this.SlaveSide.Run()
	this.LogSystem.Log("Slave Scheduler Start")
}
