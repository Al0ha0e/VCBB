package master_side

type scheduleNode struct {
	id                  string
	code                string
	baseTest            string
	hardwareRequirement string
	partitionCnt        uint64
	partitionIDOffset   uint64
	dependencies        map[string]*Dependency
	inputMap            map[string]*struct {
		x int
		y int
	}
	input           [][]string
	output          [][]string
	indeg           uint64
	outdeg          uint64
	inNodes         []*scheduleNode
	outNodes        []*scheduleNode
	controlIndeg    uint64
	controlOutdeg   uint64
	controlInNodes  []*scheduleNode
	controlOutNodes []*scheduleNode
	//minAnswerCount  uint8
}

type scheduleGraph []*scheduleNode
