package vcfs

import (
	"fmt"
	"testing"
	"time"

	"vcbb/net"

	"vcbb/peer_list"
	"vcbb/types"
)

func tEnv(ns *net.NetSimulator, addr types.Address) *FileSystem {
	pl := peer_list.NewPeerList(addr, ns)
	pl.Run()
	eg := NewMapKVStore()
	//eg := NewRedisKVStore("localhost:6379")
	return NewFileSystem(eg, pl)
}

func TestNewFileSystem(t *testing.T) {
	addr := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	ns := net.NewNetSimulator()
	fs := tEnv(ns, addr)
	fmt.Println(fs)
}

func TestSetGet(t *testing.T) {
	addr := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	ns := net.NewNetSimulator()
	fs := tEnv(ns, addr)
	err := fs.Set("test", []byte{1, 2, 3})
	err = fs.Set("test2", []byte{1, 2, 3})
	err = fs.Set("test2", []byte{2, 3, 4})
	fmt.Println(err)
	fmt.Println(fs.files["test"], fs.files["test2"])
	val, err := fs.Get("test")
	fmt.Println(val)
	val, err = fs.Get("test2")
	fmt.Println(val)
}

/*
func TestTracker(t *testing.T) {
	addr := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	ns := net.NewNetSimulator()
	pl := peer_list.NewPeerList(addr, ns)
	pl.Run()
	eg := neMapKVStore()
	fs := NewFileSystem(eg, pl)
	fs.Serve()
	pl2 := peer_list.NewPeerList(addr2, ns)
	pl2.Run()
	pl2.AddCallBack("HandleTrackerRes", func(msg peer_list.MessageInfo) {
		fmt.Println(msg, msg.From.ToString(), msg.FromSession)
		var resobj trackerRes
		json.Unmarshal(msg.Content, &resobj)
		fmt.Println(resobj)
	})
	fs.Set("test", []byte{1, 2, 3})
	fs.files["test"].peer = []types.Address{addr2}
	req := trackerReq{Keys: []string{"test"}}
	reqb, _ := json.Marshal(req)
	pl2.RemoteProcedureCall(addr, peer_list.Global, "HandleTrackerReq", reqb)
	for {
	}
}*/

func TestFileSync(t *testing.T) {
	ns := net.NewNetSimulator()
	addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	addr3 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
	fs1 := tEnv(ns, addr1)
	fs2 := tEnv(ns, addr2)
	fs3 := tEnv(ns, addr3)
	fs1.Set("test", []byte{1, 2, 3, 4, 5})
	fs1.Set("test2", []byte{2, 0, 1})
	fs1.Serve()
	fs2.Serve()
	fs3.Serve()
	fs1.SyncFile([]string{"test", "test2"}, []uint64{1, 1}, []types.Address{addr2, addr3})
	time.Sleep(10000)
	V, _ := fs2.Get("test")
	V2, _ := fs2.Get("test2")
	fmt.Println("FS2", fs2.files["test"], fs2.files["test2"])
	fmt.Println("CONT", V, V2)
	fmt.Println("FS1", fs1.files["test"], fs1.files["test2"])
	fmt.Println("FS3", fs3.files["test"], fs3.files["test2"])
	VV, _ := fs3.Get("test")
	VV2, _ := fs3.Get("test2")
	fmt.Println(VV, VV2)
}

func TestPurchaseFile(t *testing.T) {
	ns := net.NewNetSimulator()
	addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	fs1 := tEnv(ns, addr1)
	fs2 := tEnv(ns, addr2)
	fs1.Set("test", []byte{1, 2, 3, 4, 5})
	fs1.Set("test2", []byte{2, 0, 1})
	fs1.Serve()
	fs2.Serve()
	oksign := make(chan struct{}, 1)
	parts := make([]FilePart, 1)
	parts[0] = FilePart{
		Keys:  []string{"test", "test2"},
		Peers: []types.Address{addr1},
	}
	go fs2.FetchFiles(parts, oksign)
	<-oksign
	fmt.Println("FS1", fs1.files["test"], fs1.files["test2"])
	fmt.Println("FS2", fs2.files["test"], fs2.files["test2"])
	V, _ := fs2.Get("test")
	VV, _ := fs2.Get("test2")
	fmt.Println(V, VV)
}
