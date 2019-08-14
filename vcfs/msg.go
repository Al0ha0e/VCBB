package vcfs

import "vcbb/types"

type trackerReq struct {
	Keys []string `json:"keys"`
}

type trackerRes struct {
	Keys  []string          `json:"keys"`
	Peers [][]types.Address `json:"peers"`
}

type syncReq struct {
	Keys []string `json:"keys"`
	Size []uint64 `json:"size"`
}

type syncRes struct {
	Keys []string `json:"keys"`
}

type filePack struct {
	Keys  []string `json:"keys"`
	Files [][]byte `json:"files"`
}

type purchaseReq struct {
	Contract types.Address `json:"contract"`
	Keys     []string      `json:"keys"`
}

type purchaseRes struct {
	Transaction types.Address `json:"tx"`
	Keys        []string      `json:"keys"`
	Files       [][]byte      `json:"files"`
}
