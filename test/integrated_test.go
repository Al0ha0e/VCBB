package test

import (
	"fmt"
	"testing"
	"vcbb/master_side"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/slave_side"
	"vcbb/types"
	"vcbb/vcfs"
)

const url = "http://127.0.0.1:8080/execute"

func getSch(addr types.Address, ns *net.NetSimulator) *slave_side.Scheduler {
	pl := peer_list.NewPeerList(addr, ns)
	pl.Run()
	eg := vcfs.NewRedisKVStore("localhost:6379", 0)
	fs := vcfs.NewFileSystem(eg, pl)
	fs.Serve()
	exe := slave_side.NewPyExecuter(url)
	sch := slave_side.NewScheduler(100, pl, fs, exe)
	return sch
}

/*
addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
addr3 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
*/
func TestMasterAndSlave(t *testing.T) {
	addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	ns := net.NewNetSimulator()
	cli := master_side.NewMasterClient(addr1, ns)
	cli.PeerList.Peers = append(cli.PeerList.Peers, addr2)
	cli.Run()
	sch := getSch(addr2, ns)
	sch.Run()
	fmt.Println("INITOK")
	for {
	}
}
