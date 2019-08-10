package master_side

import (
	"fmt"

	"github.com/Al0ha0e/vcbb/blockchain"
	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
	"github.com/Al0ha0e/vcbb/vcfs"
)

type Scheduler struct {
	ID                 string
	bcHandler          blockchain.BlockChainHandler
	peerList           *peer_list.PeerList
	fileSystem         *vcfs.FileSystem
	graph              scheduleGraph
	result             chan *JobMeta
	originalPartitions []string
	originalData       types.DataSource
	//oriDataTransportSession *data.DataTransportSession
	//originalDataTracker     *data.Tracker
}

func NewScheduler(
	id string,
	bchandler blockchain.BlockChainHandler,
	peerList *peer_list.PeerList,
	fileSystem *vcfs.FileSystem,
	graph scheduleGraph,
	oriPartitions []string,
	oridata types.DataSource,
) (*Scheduler, error) {
	err := checkGraph(graph)
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		ID:                 id,
		bcHandler:          bchandler,
		peerList:           peerList,
		fileSystem:         fileSystem,
		graph:              graph,
		result:             make(chan *JobMeta, 100),
		originalPartitions: oriPartitions,
		originalData:       oridata,
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
	tot := this.originalData.GetTotalNum()
	for i := uint64(0); i < tot; i++ {
		sing, err := this.originalData.GetSingle(i)
		if err != nil {
			return err
		}
		err = this.fileSystem.Set(sing.GetHash(), sing.GetValue())
		if err != nil && err.Error() != "file has already setteled" {
			return err
		}
	}
	oriDep := new(Dependency)
	oriDep.DependencyJobMeta = &JobMeta{
		Id:               this.ID,
		Participants:     []types.Address{this.peerList.Address},
		Partitions:       this.originalPartitions,
		PartitionAnswers: this.originalData.GetHashList(),
	}
	for _, node := range this.graph {
		if node.indeg+node.controlIndeg == 0 {
			node.dependencies = []*Dependency{oriDep}
			err := this.DispatchJob("RANDOM ID", node)
			if err != nil {
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
			to.dependencies = append(to.dependencies, &Dependency{DependencyJob: jobmeta.job, DependencyJobMeta: jobmeta})
			if to.indeg+to.controlIndeg == 0 {
				this.DispatchJob("RANDOMID", to) //ERROR HANDLEING
			}
		}
		for _, to := range node.controlOutNodes {
			to.controlIndeg--
			if to.indeg+to.controlIndeg == 0 {
				this.DispatchJob("RANODMID", to) //ERROR HANDLEING
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
