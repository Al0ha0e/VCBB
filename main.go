package main

import (
	"fmt"
	"vcbb/types"
)

const (
	defaultDBAddr        = "localhost:6379"
	defaultBCurl         = "ws://127.0.0.1:8546"
	defaultMasterClient  = ":8080"
	defaultSlaveExecuter = "http://127.0.0.1:8080"
)

func main() {
	var account, privateKey, inaddr, outaddr, dbaddr, bcurl string
	var db int
	fmt.Println("INPUT ACCOUNT AND PRIVATEKEY")
	fmt.Scanln(&account, &privateKey)
	fmt.Println("INPUT INADDR")
	fmt.Scanln(&inaddr)
	fmt.Println("INPUT OUTADDR")
	fmt.Scanln(&outaddr)
	fmt.Println("INPUT DBADDR AND DB")
	fmt.Scanln(&dbaddr, &db)
	fmt.Println("INPUT BCURL")
	fmt.Scanln(&bcurl)
	if dbaddr == "default" {
		dbaddr = defaultDBAddr
	}
	if bcurl == "default" {
		bcurl = defaultBCurl
	}
	application, err := NewApplication(account, privateKey, inaddr, outaddr, dbaddr, bcurl, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	application.Run()
	var mode string
	fmt.Println("INPUT MODE")
	fmt.Scanln(&mode)
	if mode == "master" {
		var client string
		fmt.Println("INPUT CLIENT URL")
		fmt.Scanln(&client)
		if client == "default" {
			client = defaultMasterClient
		}
		err := application.StartMasterClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		var executer string
		fmt.Println("INPUT EXECUTER URL")
		fmt.Scanln(&executer)
		if executer == "default" {
			executer = defaultSlaveExecuter
		}
		application.StartSlaveSide(executer)
	}
	for {
		var instruct string
		fmt.Scanln(&instruct)
		if instruct == "addPeer" {
			var account, udpAddr string
			fmt.Println("INPUT ACCOUNT")
			fmt.Scanln(&account)
			fmt.Println("INPUT UDP ADDR")
			fmt.Scanln(&udpAddr)
			application.NetService.AddPeer(types.NewAddress(account), udpAddr)
			application.PeerList.Peers = append(application.PeerList.Peers, types.NewAddress(account))
		}
	}
}
