package data

import "github.com/Al0ha0e/vcbb/types"

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

type dataInfoRes struct {
	DataReceivers map[string][]types.Address `json:"dataReceivers"`
}
