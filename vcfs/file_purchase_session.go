package vcfs

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

const waitTime int64 = 10000

type filePurchaseResult struct {
	key     string
	success bool
}

type filePart struct {
	keys  []string
	peers []types.Address
}

type FilePurchaseSession struct {
	id         string
	fileSystem *FileSystem
	peerList   *peer_list.PeerListInstance
	keyMap     map[string]uint8
	parts      []filePart
	peers      []map[string]uint8
	peerChan   []chan types.Address
	stopSignal []chan struct{}
	resultChan chan filePurchaseResult
	lock       sync.Mutex
	rwlock     sync.RWMutex
	contract   types.Address
}

func NewFilePurchaseSession(id string, fs *FileSystem, parts []filePart, resultChan chan filePurchaseResult) *FilePurchaseSession {
	ret := &FilePurchaseSession{
		id:         id,
		fileSystem: fs,
		parts:      parts,
		keyMap:     make(map[string]uint8),
		resultChan: resultChan,
	}
	var i uint8
	for _, pt := range parts {
		for _, key := range pt.keys {
			ret.keyMap[key] = i
			i++
		}
	}
	ret.peers = make([]map[string]uint8, i)
	ret.peerChan = make([]chan types.Address, i)
	ret.stopSignal = make([]chan struct{}, i)
	return ret
}

func (this *FilePurchaseSession) StartSession() {
	l := len(this.peers)
	for i := 0; i < l; i++ {
		this.peers[i] = make(map[string]uint8)
		this.peerChan[i] = make(chan types.Address, 5)
		this.stopSignal[i] = make(chan struct{}, 1)
	}
	//TODO START CONTRACT
}

func (this *FilePurchaseSession) HandleTrackerRes(req peer_list.MessageInfo) {
	var reqobj trackerRes
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	for i, key := range reqobj.Keys {
		info := this.fileSystem.files[key]
		info.rwlock.RLock()
		if info.state == fPossess {
			info.rwlock.RUnlock()
			continue
		}
		this.lock.Lock()
		id := this.keyMap[key]
		for _, peer := range reqobj.Peers[i] {
			if this.peers[id][peer.ToString()] != 1 {
				this.peers[id][peer.ToString()] = 1
				this.peerChan[id] <- peer
			}
		}
		this.lock.Unlock()
		info.rwlock.RUnlock()
	}
}

func (this *FilePurchaseSession) HandlePurchaseRes(res peer_list.MessageInfo) {
	var resobj purchaseRes
	err := json.Unmarshal(res.Content, &resobj)
	if err != nil {
		return
	}
	for i, key := range resobj.Keys {
		id := this.keyMap[key]
		if this.peers[id][res.From.ToString()] != 1 {
			continue
		}
		file := resobj.Files[i]
		//TODO: CHECK FILE
		if file == nil {
			continue
		}
		info := this.fileSystem.files[key]
		info.lock.Lock()
		if info.state == fPossess {
			info.lock.Unlock()
			continue
		}
		err := this.fileSystem.engine.Set(key, file)
		if err != nil {
			info.lock.Unlock()
			continue
		}
		info.state = fPossess
		this.resultChan <- filePurchaseResult{
			key:     key,
			success: true,
		}
		close(this.peerChan[id])
		var sign struct{}
		this.stopSignal[id] <- sign
		close(this.stopSignal[id])
		info.lock.Unlock()
	}
}

func (this *FilePurchaseSession) tryToPurchase(key string) {
	id := this.keyMap[key]
	req := purchaseReq{
		Contract: this.contract,
		Keys:     []string{key},
	}
	reqb, _ := json.Marshal(req)
	for {
		select {
		case peer := <-this.peerChan[id]:
			info := this.fileSystem.files[key]
			info.rwlock.RLock()
			if info.state == fPossess {
				info.rwlock.RUnlock()
				return
			}
			info.rwlock.RUnlock()
			this.peerList.GlobalRemoteProcedureCall(peer, "HandleFilePurchaseReq", reqb)
			time.Sleep(time.Duration(waitTime))
		case <-this.stopSignal[id]:
			return
		}
	}
}
