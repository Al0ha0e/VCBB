package vcfs

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

func (this *FileSystem) HandleTrackerReq(req peer_list.MessageInfo) {
	var reqobj trackerReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	retpeer := make([][]types.Address, len(reqobj.Keys))
	for i, id := range reqobj.Keys {
		this.rwlock.RLock()
		info := this.files[id]
		this.rwlock.RUnlock()
		if info == nil {
			retpeer[i] = make([]types.Address, 0)
		} else {
			retpeer[i] = info.peer
			if info.state == fPossess {
				retpeer[i] = append(retpeer[i], this.peerList.Address)
			}
		}
	}
	res := trackerRes{
		Keys:  reqobj.Keys,
		Peers: retpeer,
	}
	resb, err := json.Marshal(res)
	if err != nil {
		return
	}
	this.peerList.Reply(req, "HandleTrackerRes", resb)
}
