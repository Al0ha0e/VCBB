package vcfs

import "github.com/Al0ha0e/vcbb/types"

type trackerReq struct {
	Keys []string `json:"keys"`
}

type trackerRes struct {
	Keys  []string          `json:"keys"`
	Peers [][]types.Address `json:"peers"`
}
