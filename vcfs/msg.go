package vcfs

import "github.com/Al0ha0e/vcbb/types"

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
