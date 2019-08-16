package vcfs

import (
	"fmt"
	"sync"

	"vcbb/peer_list"
	"vcbb/types"
)

type peerState uint8

const (
	unkown peerState = iota
	possess
	waiting
	sending
)

type fileState uint8

const (
	fUnkown fileState = iota
	fPossess
	fWaiting
	fPurchasing
	fToPurchase
)

type fileInfo struct {
	id     string
	state  fileState //TODO: CHECK STATE
	peer   []types.Address
	ps     map[string]peerState //TODO: CHECK STATE
	lock   sync.Mutex
	rwlock sync.RWMutex
}

func NewFileInfo(id string, state fileState) *fileInfo {
	return &fileInfo{
		id:    id,
		state: state,
		peer:  make([]types.Address, 0, 5),
		ps:    make(map[string]peerState),
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
	this.peerList.AddCallBack("HandleSyncReq", this.HandleSyncReq)
	this.peerList.AddCallBack("HandleSyncRes", this.HandleSyncRes)
	this.peerList.AddCallBack("HandleSyncFileArrive", this.HandleSyncFileArrive)
	this.peerList.AddCallBack("HandleFilePurchaseReq", this.HandleFilePurchaseReq)
}

func (this *FileSystem) SetInfo(key string) error {
	this.lock.Lock()
	info := this.files[key]
	if info != nil {
		this.lock.Unlock()
		return fmt.Errorf("file has already setteled")
	}
	info = NewFileInfo(key, fPossess)
	this.files[key] = info
	this.lock.Unlock()
	return nil
}

func (this *FileSystem) Set(key string, value []byte) error {
	this.lock.Lock()
	info := this.files[key]
	if info != nil {
		this.lock.Unlock()
		return fmt.Errorf("file has already setteled")
	}
	info = NewFileInfo(key, fUnkown)
	this.files[key] = info
	this.lock.Unlock()
	info.lock.Lock()
	err := this.engine.Set(key, value)
	if err != nil {
		this.lock.Lock()
		this.files[key] = nil
		this.lock.Unlock()
		info.lock.Unlock()
		return err
	}
	info.state = fPossess

	info.lock.Unlock()
	return nil
}

func (this *FileSystem) Get(key string) ([]byte, error) {
	return this.engine.Get(key)
}
