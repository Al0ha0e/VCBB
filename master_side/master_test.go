package master_side

import (
	"testing"
	"vcbb/net"
	"vcbb/types"
)

/*
addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
addr3 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
*/
/*
func TestNewScheduler(t *testing.T) {
	addr := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	ns := net.NewNetSimulator()
	pl := peer_list.NewPeerList(addr, ns)
	fs := vcfs.NewFileSystem(vcfs.NewMapKVStore(), pl)
	sch, _ := NewScheduler("TEST SCH", nil, pl, fs, nil, nil, nil)
	fmt.Println(sch)
}*/

/*
func TestNewJob(t *testing.T) {
	job := NewJob("TEST",)
}*/

func TestClient(t *testing.T) {
	//addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr1 := types.NewAddress("0x2acaac851b020ceb644bc506a3a932f4d0867afd")
	pvi := "3b82b9641714c4bb9a3e3a23ca9e8170772fcdeedd9e4591e7d03ebe564a579e"
	acco := types.NewAccount(addr1, pvi)
	ns := net.NewNetSimulator()
	cli, _ := NewMasterClient("http://127.0.0.1:8545", acco, ns)
	cli.Run()
}
