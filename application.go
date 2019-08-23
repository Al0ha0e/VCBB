package main

import (
	"vcbb/blockchain"
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
}

func NewApplication(account, privateKey, inaddr, outaddr, dbaddr, bcurl string, db int) (ret *Application, err error) {
	ret = &Application{}
	acco := types.NewAddress(account)
	ret.Account = types.NewAccount(acco, privateKey)
	ret.NetService, err = net.NewUDPNetService(inaddr, outaddr)
	if err != nil {
		return nil, err
	}
	ret.PeerList = peer_list.NewPeerList(acco, ret.NetService)
	engine, err := vcfs.NewRedisKVStore(dbaddr, db)
	if err != nil {
		return nil, err
	}
	ret.FileSystem = vcfs.NewFileSystem(engine, ret.PeerList)
	ret.BlockchainHandler, err = blockchain.NewEthBlockChainHandler(bcurl, ret.Account)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (this *Application) Run() {
	this.FileSystem.Serve()
	this.NetService.Run()
	this.PeerList.Run()
}

func (this *Application) StartMasterClient(client string) error {
	var err error
	this.MasterClient, err = master_side.NewMasterClient(this.PeerList, this.FileSystem, this.BlockchainHandler)
	if err != nil {
		return err
	}
	this.MasterClient.Run(client)
	return nil
}

func (this *Application) StartSlaveSide(executer string) {
	exe := slave_side.NewPyExecuter(executer)
	this.SlaveSide = slave_side.NewScheduler(10, this.PeerList, this.FileSystem, this.BlockchainHandler, exe)
	this.SlaveSide.Run()
}
