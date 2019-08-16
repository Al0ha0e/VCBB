package master_side

type pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type rawScheduleNode struct {
	ID                  string              `json:"id"`
	Code                string              `json:"code"`
	BaseTest            string              `json:"baseTest"`
	HardwareRequirement string              `json:"hardwareRequirement"`
	PartitionCnt        uint64              `json:"partitionCnt"`
	PartitionIDOffset   uint64              `json:"partitionIDOffset"`
	Dependencies        map[string][]string `json:"dependencies"`
	InputMap            map[string]pos      `json:"inputMap"`
	Output              [][]string          `json:"output"`
	Indeg               uint8               `json:"indeg"`
	Outdeg              uint8               `json:"outdeg"`
	OutNodes            []string            `json:"outNodes"`
	MinAnswerCount      uint8               `json:"minAnswerCount"`
}

type schReq struct {
	OriDataHash map[string]string `json:"oriDataHash"`
	SchGraph    []rawScheduleNode `json:"schGraph"`
}
