package data

type dataTransportReq struct {
	Requirement string   `json:"requirement"`
	Metadata    []string `json:"meta"`
}

type dataTransportRes struct {
	Amount uint64 `json:"amount"`
}

type dataReceivedRes struct {
	//Success  bool   `json:"success"`
	DataList []uint64 `json:"dataList"`
}
