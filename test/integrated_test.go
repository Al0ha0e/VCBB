package test

import (
	"fmt"
	"testing"
	"vcbb/blockchain"
	"vcbb/master_side"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/slave_side"
	"vcbb/types"
	"vcbb/vcfs"
)

const url = "localhost:8080"

func getSch(account *types.Account, ns *net.NetSimulator) *slave_side.Scheduler {
	pl := peer_list.NewPeerList(account.Id, ns)
	pl.Run()
	eg, _ := vcfs.NewRedisKVStore("localhost:6379", 0)
	fs := vcfs.NewFileSystem(eg, pl)
	fs.Serve()
	exe := slave_side.NewPyExecuter(url)
	bch, _ := blockchain.NewEthBlockChainHandler("ws://127.0.0.1:8546", account)
	sch := slave_side.NewScheduler(100, pl, fs, bch, exe)
	return sch
}

/*
addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
addr3 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
*/
func TestMasterAndSlave(t *testing.T) {
	addr1 := types.NewAddress("0x53143B153321E9cBC268C806aFeF444346AA6522" /*"0x2acaac851b020ceb644bc506a3a932f4d0867afd"*/)
	pri1 := "52533c066a5b3ddc6848fd692b93abc6025d06aa5a220f67213c619a6ded8d3a" //"3b82b9641714c4bb9a3e3a23ca9e8170772fcdeedd9e4591e7d03ebe564a579e"
	addr2 := types.NewAddress("0xe3b2d0E78a88cB0139aFC7A08733766F5C788A12" /*"0x9c67d6e615fb9fb28ddad773fbcfa8e5dad092f3"*/)
	pri2 := "680771f89fa288c71a113846de2b953bfedb6da300e1747a1d1b51f2a43a10c5" //"ee09c465edc1674d382157f9edb26681707b79b31cab452450776a2a1ad57be5"
	acco1 := types.NewAccount(addr1, pri1)
	acco2 := types.NewAccount(addr2, pri2)
	ns := net.NewNetSimulator()
	pl := peer_list.NewPeerList(addr1, ns)
	pl.Run()
	eg, _ := vcfs.NewRedisKVStore("localhost:6379", 1)
	fs := vcfs.NewFileSystem(eg, pl)
	fs.Serve()
	bc, _ := blockchain.NewEthBlockChainHandler("ws://127.0.0.1:8546", acco1)
	cli, _ := master_side.NewMasterClient(pl, fs /*"ws://127.0.0.1:8546"*/, bc /*, acco1, ns*/)
	cli.PeerList.Peers = append(cli.PeerList.Peers, addr2)
	cli.Run(":8080")
	sch := getSch(acco2, ns)
	sch.Run()
	fmt.Println("INITOK")
	for {
	}
}
