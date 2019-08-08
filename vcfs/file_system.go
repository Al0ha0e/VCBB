package vcfs

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

type fileInfo struct {
	id        string
	local     bool
	peer      []types.Address
	peerState map[string]uint8
}

type FileSystem struct {
	engine   *KVStore
	files    map[string]*fileInfo
	peerList *peer_list.PeerList
}

func NewFileSystem(eg *KVStore, pl *peer_list.PeerList) *FileSystem {
	return &FileSystem{
		engine:   eg,
		files:    make(map[string]*fileInfo),
		peerList: pl,
	}
}

func (this *FileSystem) Serve() {
	this.peerList.AddCallBack("HandleTracker", this.HandleTracker)
}

func (this *FileSystem) HandleTracker(msg peer_list.MessageInfo) {
	var msgobj trackerReq
	err := json.Unmarshal(msg.Content, &msgobj)
	if err != nil {
		return
	}
	retpeer := make([][]types.Address, len(msgobj.Keys), 0)
	for i, id := range msgobj.Keys {
		info := this.files[id]
		if info == nil {
			retpeer[i] = make([]types.Address, 0, 0)
		} else {
			retpeer[i] = info.peer
		}
	}
	res := trackerRes{
		Keys:  msgobj.Keys,
		Peers: retpeer,
	}
	resb, err := json.Marshal(res)
	if err != nil {
		return
	}
	this.peerList.RemoteProcedureCall(msg.From, "", resb)
}
