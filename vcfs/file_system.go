package vcfs

import (
	"sync"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

type fileState uint8

const (
	unkown fileState = iota
	possess
	waiting
	sending
)

type fileInfo struct {
	id        string
	local     bool
	peer      []types.Address
	peerState map[string]fileState
	lock      sync.Mutex
	rwlock    sync.RWMutex
}

func NewFileInfo(id string, local bool) *fileInfo {
	return &fileInfo{
		id:        id,
		local:     local,
		peer:      make([]types.Address, 0, 5),
		peerState: make(map[string]fileState),
	}
}

type FileSystem struct {
	engine   KVStore
	files    map[string]*fileInfo
	peerList *peer_list.PeerList
	lock     sync.Mutex
	rwlock   sync.RWMutex
}

func NewFileSystem(eg KVStore, pl *peer_list.PeerList) *FileSystem {
	return &FileSystem{
		engine:   eg,
		files:    make(map[string]*fileInfo),
		peerList: pl,
	}
}

func (this *FileSystem) Serve() {
	this.peerList.AddCallBack("HandleTrackerReq", this.HandleTrackerReq)
}
