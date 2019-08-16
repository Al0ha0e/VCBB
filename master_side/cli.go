package master_side

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

type MasterClient struct {
	url        string
	schedulers []*Scheduler
	peerList   *peer_list.PeerList
	fileSystem *vcfs.FileSystem
}

func NewMasterClient(addr types.Address, ns *net.NetSimulator) *MasterClient {
	pl := peer_list.NewPeerList(addr, ns)
	ret := &MasterClient{
		schedulers: make([]*Scheduler, 0),
		peerList:   pl,
		fileSystem: vcfs.NewFileSystem(vcfs.NewRedisKVStore("localhost:6379", 1), pl),
	}
	return ret
}

func (this *MasterClient) Run() {
	fmt.Println("ST")
	this.peerList.Run()
	this.fileSystem.Serve()
	http.HandleFunc("/commitSchGraph", this.handleReq)
	http.ListenAndServe(":8080", nil)
}

func (this *MasterClient) handleReq(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	var reqobj schReq
	err = json.Unmarshal(body, &reqobj)
	if err != nil {
		return
	}
	graph, err := this.ConstructGraph(reqobj.SchGraph)
	if err != nil {
		return
	}
	NewScheduler("test", nil, this.peerList, this.fileSystem, graph, reqobj.OriDataHash)
	/*
		fmt.Println(sb.OriDataHash)
		for _, node := range sb.SchGraph {
			fmt.Println(node)
		}*/
}

func (this *MasterClient) ConstructGraph(rawGraph []rawScheduleNode) (scheduleGraph, error) {
	ret := make([]*scheduleNode, len(rawGraph))
	idmap := make(map[string]*scheduleNode)
	for i, rawnode := range rawGraph {
		deps := make(map[string]*Dependency)
		for k, v := range rawnode.Dependencies {
			deps[k] = &Dependency{
				keys: v,
			}
		}
		inpt := make([][]string, rawnode.PartitionCnt)
		for j := 0; j < int(rawnode.PartitionCnt); j++ {
			inpt[j] = make([]string, rawnode.InputCnt)
		}
		node := NewScheduleNode(rawnode.ID, rawnode.Code, rawnode.BaseTest, rawnode.HardwareRequirement,
			rawnode.PartitionCnt, rawnode.PartitionIDOffset, deps, rawnode.InputMap, inpt, rawnode.Output, rawnode.Indeg,
			rawnode.Outdeg, rawnode.MinAnswerCount)
		ret[i] = node
		idmap[node.id] = node
	}
	for i, rawnode := range rawGraph {
		inNodes := make([]*scheduleNode, 0)
		outNodes := make([]*scheduleNode, len(rawnode.OutNodes))
		for id, _ := range rawnode.Dependencies {
			inNodes = append(inNodes, idmap[id])
		}
		for j, id := range rawnode.OutNodes {
			outNodes[j] = idmap[id]
		}
		ret[i].inNodes, ret[i].outNodes = inNodes, outNodes
	}
	return ret, nil
}
