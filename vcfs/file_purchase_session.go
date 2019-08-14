package vcfs

import (
	"encoding/json"
	"sync"
	"time"

	"vcbb/peer_list"
	"vcbb/types"
)

const (
	waitTime    int64 = 10000
	checkPeriod int64 = 50000
)

type filePurchaseResult struct {
	key     string
	success bool
}

type FilePart struct {
	Keys  []string
	Peers []types.Address
}

type FilePurchaseSession struct {
	id                string
	fileSystem        *FileSystem
	peerList          *peer_list.PeerListInstance
	keyMap            map[string]uint8
	parts             []FilePart
	partCnt           uint8
	peers             []map[string]uint8
	peerChan          []chan types.Address
	stopSignal        []chan struct{}
	rdCheckStopSignal chan struct{}
	resultChan        chan filePurchaseResult
	lock              sync.Mutex
	rwlock            sync.RWMutex
	contract          types.Address
}

func NewFilePurchaseSession(id string, fs *FileSystem, parts []FilePart, resultChan chan filePurchaseResult) *FilePurchaseSession {
	ret := &FilePurchaseSession{
		id:         id,
		fileSystem: fs,
		parts:      parts,
		partCnt:    uint8(len(parts)),
		keyMap:     make(map[string]uint8),
		resultChan: resultChan,
	}
	var i uint8
	for _, pt := range parts {
		for _, key := range pt.Keys {
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
	//TODO START CONTRACT
	for i := 0; i < l; i++ {
		this.peers[i] = make(map[string]uint8)
		this.peerChan[i] = make(chan types.Address, 5)
		this.stopSignal[i] = make(chan struct{}, 1)
	}
	this.rdCheckStopSignal = make(chan struct{})
	this.peerList = this.fileSystem.peerList.GetInstance(this.id)
	this.peerList.AddCallBack("HandleTrackerRes", this.HandleTrackerRes)
	this.peerList.AddCallBack("HandlePurchaseRes", this.HandlePurchaseRes)
	go this.RoundCheck()
	for key, value := range this.keyMap {
		go this.tryToPurchase(key, value)
	}
	for _, part := range this.parts {
		req := trackerReq{
			Keys: part.Keys,
		}
		reqb, _ := json.Marshal(req)
		for _, peer := range part.Peers {
			prstr := peer.ToString()
			for _, key := range part.Keys {
				this.peers[this.keyMap[key]][prstr] = 1
			}
			this.peerList.GlobalRemoteProcedureCall(peer, "HandleTrackerReq", reqb)
		}
	}
}

func (this *FilePurchaseSession) HandleTrackerRes(req peer_list.MessageInfo) {
	var reqobj trackerRes
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	//fmt.Println("RES", reqobj)
	for i, key := range reqobj.Keys {
		id := this.keyMap[key]
		if this.peers[id][req.From.ToString()] != 1 {
			continue
		}
		info := this.fileSystem.files[key]
		info.rwlock.RLock()
		if info.state == fPossess {
			info.rwlock.RUnlock()
			continue
		}
		this.lock.Lock()
		for _, peer := range reqobj.Peers[i] {
			if this.peers[id][peer.ToString()] != 2 {
				this.peers[id][peer.ToString()] = 2
				this.peerChan[id] <- peer
			}
		}
		this.lock.Unlock()
		info.rwlock.RUnlock()
	}
}

func (this *FilePurchaseSession) HandlePurchaseRes(res peer_list.MessageInfo) {
	var resobj purchaseRes
	fr := res.From.ToString()
	err := json.Unmarshal(res.Content, &resobj)
	if err != nil {
		return
	}
	//fmt.Println("PURCHASE RES", resobj)
	for i, key := range resobj.Keys {
		id := this.keyMap[key]
		if this.peers[id][fr] != 2 {
			continue
		}
		file := resobj.Files[i]
		//TODO: CHECK FILE
		if file == nil {
			continue
		}
		info := this.fileSystem.files[key]
		info.lock.Lock()
		if info.ps[fr] != possess {
			info.ps[fr] = possess
			info.peer = append(info.peer, res.From)
		}
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
		var sign struct{}
		this.stopSignal[id] <- sign
		close(this.stopSignal[id])
		close(this.peerChan[id])
		this.partCnt--
		if this.partCnt == 0 {
			this.Terminate()
		}
		this.resultChan <- filePurchaseResult{
			key:     key,
			success: true,
		}
		info.lock.Unlock()
	}
}

func (this *FilePurchaseSession) tryToPurchase(key string, id uint8) {
	req := purchaseReq{
		Contract: this.contract,
		Keys:     []string{key},
	}
	reqb, _ := json.Marshal(req)
	for {
		select {
		case peer := <-this.peerChan[id]:
			//fmt.Println("TRY", peer.ToString())
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

func (this *FilePurchaseSession) RoundCheck() {
	for {
		select {
		case <-this.rdCheckStopSignal:
			return
		default:
			time.Sleep(time.Duration(checkPeriod))
			for key, id := range this.keyMap {
				info := this.fileSystem.files[key]
				info.rwlock.RLock()
				if info.state == fPossess {
					_, ok := <-this.peerChan[id]
					if ok {
						var sign struct{}
						this.stopSignal[id] <- sign
						close(this.stopSignal[id])
						close(this.peerChan[id])
						this.partCnt--
						if this.partCnt == 0 {
							this.Terminate()
						}
						this.resultChan <- filePurchaseResult{
							key:     key,
							success: true,
						}
					}
				}
				info.rwlock.RUnlock()
			}
		}
	}
}

func (this *FilePurchaseSession) Terminate() {
	var sign struct{}
	this.rdCheckStopSignal <- sign
	close(this.rdCheckStopSignal)
	this.fileSystem.peerList.RemoveInstance(this.id)
	//TODO: UPDATE FILE INFO PEER STATE
}
