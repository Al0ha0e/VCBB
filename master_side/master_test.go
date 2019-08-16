package master_side

import (
	"testing"
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
	cli := new(MasterClient)
	cli.Run()
	for {
	}
}
