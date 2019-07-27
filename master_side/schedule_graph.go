package master_side

type scheduleNode struct {
	id                  string
	code                string
	baseTest            string
	hardwareRequirement string
	partitions          []string
	dependencies        []*Dependency
	indeg               uint64
	outdeg              uint64
	inNodes             []*scheduleNode
	outNodes            []*scheduleNode
	controlIndeg        uint64
	controlOutdeg       uint64
	controlInNodes      []*scheduleNode
	controlOutNodes     []*scheduleNode
}

type scheduleGraph []*scheduleNode
