package master_side

type position struct {
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
	InputCnt            uint64              `json:"inputCnt"`
	InputMap            map[string]position `json:"inputMap"`
	Output              [][]string          `json:"output"`
	Indeg               uint64              `json:"indeg"`
	Outdeg              uint64              `json:"outdeg"`
	OutNodes            []string            `json:"outNodes"`
	MinAnswerCount      uint8               `json:"minAnswerCount"`
}

type schReq struct {
	OriDataHash map[string]string `json:"oriDataHash"`
	SchGraph    []rawScheduleNode `json:"schGraph"`
}
