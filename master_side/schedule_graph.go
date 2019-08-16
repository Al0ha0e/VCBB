package master_side

type scheduleNode struct {
	id                  string
	code                string
	baseTest            string
	hardwareRequirement string
	partitionCnt        uint64
	partitionIDOffset   uint64
	dependencies        map[string]*Dependency
	inputMap            map[string]position
	input               [][]string
	output              [][]string
	indeg               uint64
	outdeg              uint64
	inNodes             []*scheduleNode
	outNodes            []*scheduleNode
	controlIndeg        uint64
	controlOutdeg       uint64
	controlInNodes      []*scheduleNode
	controlOutNodes     []*scheduleNode
	minAnswerCount      uint8
}

type scheduleGraph []*scheduleNode

func NewScheduleNode(
	id,
	code,
	baseTest,
	hardwareRequirement string,
	partitionCnt,
	partitionIDOffset uint64,
	dependencies map[string]*Dependency,
	inputMap map[string]position,
	input,
	output [][]string,
	indeg,
	outdeg uint64,
	//inNodes,
	//outNodes []*scheduleNode,
	/*
		controlIndeg,
		controlOutdeg uint64,
		controlInNodes,
		controlOutNodes []*scheduleNode,*/
	minAnswerCount uint8) *scheduleNode {
	return &scheduleNode{
		id:                  id,
		code:                code,
		baseTest:            baseTest,
		hardwareRequirement: hardwareRequirement,
		partitionCnt:        partitionCnt,
		partitionIDOffset:   partitionIDOffset,
		dependencies:        dependencies,
		inputMap:            inputMap,
		input:               input,
		output:              output,
		indeg:               indeg,
		outdeg:              outdeg,
		//inNodes:         inNodes,
		//outNodes:        outNodes,
		/*
			controlIndeg:    controlIndeg,
			controlOutdeg:   controlOutdeg,
			controlInNodes:  controlInNodes,
			controlOutNodes: controlOutNodes,*/
		minAnswerCount: minAnswerCount,
	}
}
