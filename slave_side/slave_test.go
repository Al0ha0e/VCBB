package slave_side

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"vcbb/msg"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

const url = "http://127.0.0.1:8080/execute"

/*
func TestPyExecuter(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cli.Set("A", "1", 0)
	cli.Set("B", "2", 0)
	exe := NewPyExecuter(url)
	code :=
		`def func():
	print(input[0],input[1],input[0]+input[1])
	return input[0]+input[1]
output=[func()]`
	ansChan := make(chan *executeResult, 1)
	exe.Run(1, [][]string{[]string{"A", "B"}}, code, ansChan)
	ans := <-ansChan
	if ans.err != nil {
		fmt.Println(ans.err)
		return
	}
	ansstr, _ := cli.Get(ans.result[0][0]).Result()
	println("ANS", ansstr)
}*/

/*
addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
addr3 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d3")
*/

func getSch(addr types.Address, ns *net.NetSimulator) *Scheduler {
	pl := peer_list.NewPeerList(addr, ns)
	pl.Run()
	eg := vcfs.NewRedisKVStore("localhost:6379", 0)
	fs := vcfs.NewFileSystem(eg, pl)
	fs.Serve()
	exe := NewPyExecuter(url)
	sch := NewScheduler(100, pl, fs, nil, exe)
	return sch
}

func TestHandleSeekParticipantReq(t *testing.T) {
	addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	ns := net.NewNetSimulator()
	sch := getSch(addr1, ns)
	sch.Run()
	pl2 := peer_list.NewPeerList(addr2, ns)
	pl2.Run()
	inst := pl2.GetInstance("TEST")
	inst.AddCallBack("handleMetaDataReq", func(req peer_list.MessageInfo) {
		fmt.Println(req.From.ToString(), req.FromSession, req.Content)
		var reqobj msg.MetaDataReq
		json.Unmarshal(req.Content, &reqobj)
		fmt.Println(reqobj)
	})
	req := msg.ComputationReq{
		Id:           "TEST",
		Master:       addr1,
		ContractAddr: addr2,
		PartitionCnt: 10,
		Hardware:     "TESTH",
		BaseTest:     "TESTB",
	}
	reqb, _ := json.Marshal(req)
	inst.GlobalRemoteProcedureCall(addr1, "handleSeekParticipantReq", reqb)
	time.Sleep(10000)
}

func TestHandleMetaDataRes(t *testing.T) {
	addr1 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1")
	addr2 := types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2")
	ns := net.NewNetSimulator()
	sch := getSch(addr1, ns)
	sch.Run()
	pl2 := peer_list.NewPeerList(addr2, ns)
	pl2.Run()
	fs2 := vcfs.NewFileSystem(vcfs.NewMapKVStore(), pl2)
	fs2.Serve()
	fs2.Set("A", []byte("1"))
	fs2.Set("B", []byte("2"))
	inst := pl2.GetInstance("TEST")
	inst.AddCallBack("handleMetaDataReq", func(req peer_list.MessageInfo) {
		meta := msg.JobMeta{
			Participants: []types.Address{addr2},
			Keys:         []string{"A", "B"},
		}
		res := msg.MetaDataRes{
			PartitionIdOffset: 3,
			Inputs:            [][]string{[]string{"A", "B"}},
			DependencyMeta:    []msg.JobMeta{meta},
			Code: `def func():
	print(input[0],input[1],input[0]+input[1])
	print(partitionId)
	return input[0]+input[1]
output=[func()]`,
		}
		fmt.Println("RES", res)
		resb, _ := json.Marshal(res)
		inst.Reply(req, "handleMetaDataRes", resb)
	})
	req := msg.ComputationReq{
		Id:           "TEST",
		Master:       addr2,
		ContractAddr: addr2,
		PartitionCnt: 1,
		Hardware:     "TESTH",
		BaseTest:     "TESTB",
	}
	reqb, _ := json.Marshal(req)
	inst.GlobalRemoteProcedureCall(addr1, "handleSeekParticipantReq", reqb)
	for {
	}
}
