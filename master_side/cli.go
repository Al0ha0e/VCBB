package master_side

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vcbb/blockchain"
	"vcbb/log"
	"vcbb/peer_list"
	"vcbb/vcfs"
)

type MasterClient struct {
	url        string
	schedulers []*Scheduler
	PeerList   *peer_list.PeerList
	fileSystem *vcfs.FileSystem
	bcHandler  *blockchain.EthBlockChainHandler
	logger     *log.LoggerInstance
}

func NewMasterClient(pl *peer_list.PeerList, fs *vcfs.FileSystem, bchandler *blockchain.EthBlockChainHandler, logSystem *log.LogSystem) (*MasterClient, error) {
	ret := &MasterClient{
		schedulers: make([]*Scheduler, 0),
		PeerList:   pl,
		fileSystem: fs,
		bcHandler:  bchandler,
		logger:     logSystem.GetInstance(log.Topic("Master Client")),
	}
	ret.logger.Log("New Master Client")
	return ret, nil
}

func (this *MasterClient) Run(url string) {
	this.logger.Log("Start")
	//this.PeerList.Run()
	//this.fileSystem.Serve()
	http.HandleFunc("/commitSchGraph", this.handleReq)
	go http.ListenAndServe(url /*":8080"*/, nil)
}

func (this *MasterClient) handleReq(w http.ResponseWriter, req *http.Request) {
	this.logger.Log("Receive Request")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		this.logger.Err("Fail To Read Request Body " + err.Error())
		return
	}
	var reqobj schReq
	err = json.Unmarshal(body, &reqobj)
	if err != nil {
		this.logger.Err("Fail To Unmarshal Request " + err.Error())
		return
	}
	graph, err := this.constructGraph(reqobj.SchGraph)
	if err != nil {
		this.logger.Err("Fail To Construct Graph " + err.Error())
		return
	}
	result := make(chan [][]string, 1)
	sch, err := NewScheduler("test", this.bcHandler, this.PeerList, this.fileSystem, graph, reqobj.OriDataHash, result, this.logger)
	if err != nil {
		this.logger.Err("Fail To Create Scheduler " + err.Error())
		return
	}
	//fmt.Println(sch)
	this.logger.Log("Scheduler Create OK")
	sch.Dispatch()
	ans := <-result
	this.logger.Log("Compute Over")
	for _, partition := range ans {
		for _, anshash := range partition {
			ansdata, _ := this.fileSystem.Get(anshash)
			this.logger.Log("Final Answer Data " + string(ansdata))
		}
	}
	/*
		fmt.Println(sb.OriDataHash)
		for _, node := range sb.SchGraph {
			fmt.Println(node)
		}*/
}

func (this *MasterClient) constructGraph(rawGraph []rawScheduleNode) (scheduleGraph, error) {
	this.logger.Log("Try To Create Graph")
	ret := make([]*scheduleNode, len(rawGraph))
	idmap := make(map[string]*scheduleNode)
	this.logger.Log("Try To Construct Nodes")
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
	this.logger.Log("Try To Construct Node Relation")
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
