package vcfs

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
)

type fileInfo struct {
	id    string
	local bool
}

type FileSystem struct {
	engine      *KVStore
	files       map[string]*fileInfo
	peerList    *peer_list.PeerList
	stopSignal  chan struct{}
	trackerChan chan peer_list.MessageInfo
}

func NewFileSystem(eg *KVStore, pl *peer_list.PeerList) *FileSystem {
	return &FileSystem{
		engine:     eg,
		files:      make(map[string]*fileInfo),
		peerList:   pl,
		stopSignal: make(chan struct{}),
	}
}

func (this *FileSystem) Serve() {
	this.peerList.AddChannel("HandleTracker", this.trackerChan)
	go this.HandleTracker()
}

func (this *FileSystem) HandleTracker() {
	for {
		select {
		case <-this.stopSignal:
			return
		case msg := <-this.trackerChan:
			var msgobj trackerReq
			err := json.Unmarshal(msg.Content, &msgobj)
			if err != nil {
				continue
			}

		}
	}
}
