package master_side

import (
	"fmt"

	"vcbb/blockchain"
	"vcbb/peer_list"
	"vcbb/types"
	"vcbb/vcfs"
)

type Scheduler struct {
	ID           string
	bcHandler    blockchain.BlockChainHandler
	peerList     *peer_list.PeerList
	fileSystem   *vcfs.FileSystem
	graph        scheduleGraph
	result       chan *JobMeta
	originalData map[string]string
	idSource     *types.UniqueRandomIDSource
	//oriDataTransportSession *data.DataTransportSession
	//originalDataTracker     *data.Tracker
}

func NewScheduler(
	id string,
	bchandler blockchain.BlockChainHandler,
	peerList *peer_list.PeerList,
	fileSystem *vcfs.FileSystem,
	graph scheduleGraph,
	oridata map[string]string,
) (*Scheduler, error) {
	err := checkGraph(graph)
	if err != nil {
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
	}, nil
}

func (this *Scheduler) DispatchJob(id string, node *scheduleNode) error {
	job := NewJob(
		id, //TODO: RANDOMID
		this,
		node,
		node.minAnswerCount,
		/*
			[]*Dependency{oriDep},
			node.partitions,
			node.code,
			node.baseTest,
			node.hardwareRequirement,*/
	)
	err := job.StartSession(this) //BUG:BLOCK WHEN A JOB WAITING FOR ITS CONTRACT
	if err != nil {
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
	oriMeta := &JobMeta{
		Participants: []types.Address{this.peerList.Address},
	}
	for _, v := range this.originalData {
		this.fileSystem.SetInfo(v)
	}
	fmt.Println("LL", len(this.graph))
	for _, node := range this.graph {
		fmt.Println("TRY ID", node.id)
		if node.indeg+node.controlIndeg == 0 {
			fmt.Println("0 INDEG", node.id)
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
	for {
		jobmeta := <-this.result
		node := jobmeta.node
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
		for _, to := range node.controlOutNodes {
			to.controlIndeg--
			if to.indeg+to.controlIndeg == 0 {
				this.DispatchJob(this.idSource.Get(), to) //ERROR HANDLEING
			}
		}
	}
}

func checkGraph(graph scheduleGraph) error {
	l := len(graph)
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
