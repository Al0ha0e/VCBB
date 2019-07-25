package master_side

import (
	"fmt"
	"vcbb/blockchain"
	"vcbb/peer_list"
)

type Scheduler struct {
	bcHandler blockchain.BlockChainHandler
	peerList  *peer_list.PeerList
	graph     scheduleGraph
}

func NewScheduler(
	bchandler blockchain.BlockChainHandler,
	peerList *peer_list.PeerList,
	graph scheduleGraph) (*Scheduler, error) {
	err := checkGraph(graph)
	if err != nil {
		return nil, err
	}
	return &Scheduler{bcHandler: bchandler, peerList: peerList, graph: graph}, nil
}

func (this *Scheduler) Dispatch() {
	for _, node := range this.graph {
		if node.indeg+node.controlIndeg == 0 {

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
