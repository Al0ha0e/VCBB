package master_side

import (
	"fmt"
	"strconv"

	"vcbb/blockchain"
	"vcbb/log"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

type Scheduler struct {
	ID           string
	bcHandler    *blockchain.EthBlockChainHandler
	peerList     *peer_list.PeerList
	fileSystem   *vcfs.FileSystem
	graph        scheduleGraph
	result       chan *JobMeta
	originalData map[string]string
	idSource     *types.UniqueRandomIDSource
	finalResult  chan [][]string
	logger       *log.LoggerInstance
	//oriDataTransportSession *data.DataTransportSession
	//originalDataTracker     *data.Tracker
}

func NewScheduler(
	id string,
	bchandler *blockchain.EthBlockChainHandler,
	peerList *peer_list.PeerList,
	fileSystem *vcfs.FileSystem,
	graph scheduleGraph,
	oridata map[string]string,
	result chan [][]string,
	fatherLogger *log.LoggerInstance,
) (*Scheduler, error) {
	fatherLogger.Log("Try To Create Scheduler " + id)
	logger := fatherLogger.GetSubInstance(log.Topic(id))
	err := checkGraph(graph, logger)
	if err != nil {
		logger.Err("Fail To Check Graph " + err.Error())
		return nil, err
	}
	return &Scheduler{
		ID:           id,
		bcHandler:    bchandler,
		peerList:     peerList,
		fileSystem:   fileSystem,
		graph:        graph,
		result:       make(chan *JobMeta, 100),
		originalData: oridata,
		idSource:     types.NewUniqueRandomIDSource(32),
		finalResult:  result,
		logger:       logger,
	}, nil
}

func (this *Scheduler) DispatchJob(id string, node *scheduleNode) error {
	this.logger.Log("Try To Dispatch Job " + id)
	job := NewJob(
		id, //TODO: RANDOMID
		this,
		node,
		node.minAnswerCount,
		this.logger,
		/*
			[]*Dependency{oriDep},
			node.partitions,
			node.code,
			node.baseTest,
			node.hardwareRequirement,*/
	)
	err := job.StartSession(this) //BUG:BLOCK WHEN A JOB WAITING FOR ITS CONTRACT
	if err != nil {
		this.logger.Err("Fail To Dispatch Job " + err.Error())
		return err
	}
	return nil
}

/*
* distribute data
* start file tracker
* start all of the jobs whose indeg is 0
* watch the running state
 */
func (this *Scheduler) Dispatch() error {
	this.logger.Log("Scheduler Start")
	oriMeta := &JobMeta{
		Participants: []types.Address{this.peerList.Address},
	}
	this.logger.Log("Try To Set Original Data Info")
	for _, v := range this.originalData {
		this.fileSystem.SetInfo(v)
	}
	this.logger.Log("Searching For Zero-Indeg-Node")
	for _, node := range this.graph {
		if node.indeg+node.controlIndeg == 0 {
			node.dependencies["ori"].dependencyJobMeta = oriMeta
			for k, v := range this.originalData {
				pos, ok := node.inputMap[k]
				if ok {
					node.input[pos.X][pos.Y] = v
				}
			}
			err := this.DispatchJob(this.idSource.Get(), node)
			if err != nil {
				fmt.Println("ERROR", err)
				return err
			}
		}
	}
	go this.watch()
	return nil
}

func (this *Scheduler) watch() {
	this.logger.Log("Start Watching")
	for {
		jobmeta := <-this.result
		this.logger.Log("Job Terminated Id: " + jobmeta.Id + " Contract: " + jobmeta.Contract.ToString())
		node := jobmeta.node
		if node.outdeg+node.controlOutdeg == 0 {
			this.logger.Log("Final Job Terminated")
			var keys []string
			for _, partitions := range jobmeta.PartitionAnswers {
				for _, answer := range partitions {
					keys = append(keys, answer)
				}
			}
			part := vcfs.FilePart{
				Peers: jobmeta.Participants,
				Keys:  keys,
			}
			ok := make(chan struct{}, 1)
			this.logger.Log("Try To Fetch Final Answer")
			this.fileSystem.FetchFiles([]vcfs.FilePart{part}, ok)
			<-ok
			this.logger.Log("Final Answer Fetch OK")
			this.finalResult <- jobmeta.PartitionAnswers
			return
		}
		this.logger.Log("Update OutNodes")
		for _, to := range node.outNodes {
			to.indeg--
			//to.dependencies = append(to.dependencies, &Dependency{DependencyJob: jobmeta.job, DependencyJobMeta: jobmeta})
			for i := 0; i < len(jobmeta.PartitionAnswers); i++ {
				for j := 0; j < len(jobmeta.PartitionAnswers[i]); j++ {
					id := node.output[i][j]
					pos, ok := to.inputMap[id]
					if ok {
						to.input[pos.X][pos.Y] = jobmeta.PartitionAnswers[i][j]
					}
				}
			}
			to.dependencies[node.id].dependencyJobMeta = jobmeta
			if to.indeg+to.controlIndeg == 0 {
				this.DispatchJob(this.idSource.Get(), to) //ERROR HANDLEING
			}
		}
		this.logger.Log("Update Control OutNodes")
		for _, to := range node.controlOutNodes {
			to.controlIndeg--
			if to.indeg+to.controlIndeg == 0 {
				this.DispatchJob(this.idSource.Get(), to) //ERROR HANDLEING
			}
		}
	}
}

func checkGraph(graph scheduleGraph, logger *log.LoggerInstance) error {
	l := len(graph)
	logger.Log("Check Graph of length " + strconv.Itoa(int(l)))
	if l == 0 {
		return fmt.Errorf("empty graph")
	}
	indeg := make(map[string]uint64)
	var stack []*scheduleNode
	var outnodecnt uint64 = 0
	for i := 0; i < l; i++ {
		if graph[i].outdeg+graph[i].controlOutdeg == 0 {
			outnodecnt++
			if outnodecnt > 1 {
				return fmt.Errorf("output node more than one")
			}
		}
		indeg[graph[i].id] = graph[i].indeg + graph[i].controlIndeg
		if graph[i].indeg+graph[i].controlIndeg == 0 {
			stack = append(stack, graph[i])
		}
	}
	if outnodecnt == 0 {
		return fmt.Errorf("graph is not a DAG")
	}
	vis := make(map[string]bool)
	for len(stack) > 0 {
		obj := stack[0]
		vis[obj.id] = true
		stack = stack[1:]
		for _, o := range obj.outNodes {
			id := o.id
			if vis[id] {
				return fmt.Errorf("graph is not a DAG")
			}
			indeg[id]--
			if indeg[id] == 0 {
				stack = append(stack, obj)
			}
		}
		for _, o := range obj.controlOutNodes {
			id := o.id
			if vis[id] {
				return fmt.Errorf("graph is not a DAG")
			}
			indeg[id]--
			if indeg[id] == 0 {
				stack = append(stack, obj)
			}
		}
	}
	return nil
}
