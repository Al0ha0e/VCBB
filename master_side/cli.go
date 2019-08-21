package master_side

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vcbb/blockchain"
	"vcbb/net"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

type MasterClient struct {
	url        string
	schedulers []*Scheduler
	PeerList   *peer_list.PeerList
	fileSystem *vcfs.FileSystem
	bcHandler  *blockchain.EthBlockChainHandler
}

func NewMasterClient(bcurl string, account *types.Account, ns *net.NetSimulator) (*MasterClient, error) {
	pl := peer_list.NewPeerList(account.Id, ns)
	bchandler, err := blockchain.NewEthBlockChainHandler(bcurl, account)
	if err != nil {
		return nil, err
	}
	kv, err := vcfs.NewRedisKVStore("localhost:6379", 1)
	if err != nil {
		return nil, err
	}
	ret := &MasterClient{
		schedulers: make([]*Scheduler, 0),
		PeerList:   pl,
		fileSystem: vcfs.NewFileSystem(kv, pl),
		bcHandler:  bchandler,
	}
	return ret, nil
}

func (this *MasterClient) Run() {
	fmt.Println("ST")
	this.PeerList.Run()
	this.fileSystem.Serve()
	http.HandleFunc("/commitSchGraph", this.handleReq)
	go http.ListenAndServe(":8080", nil)
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
	graph, err := this.constructGraph(reqobj.SchGraph)
	if err != nil {
		return
	}
	sch, err := NewScheduler("test", this.bcHandler, this.PeerList, this.fileSystem, graph, reqobj.OriDataHash)
	if err != nil {
		return
	}
	//fmt.Println(sch)
	sch.Dispatch()
	/*
		fmt.Println(sb.OriDataHash)
		for _, node := range sb.SchGraph {
			fmt.Println(node)
		}*/
}

func (this *MasterClient) constructGraph(rawGraph []rawScheduleNode) (scheduleGraph, error) {
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
