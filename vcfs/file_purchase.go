package vcfs

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"vcbb/peer_list"
)

const waitRoundCheckDuration int64 = 10000

func (this *FileSystem) FetchFiles(parts []FilePart, okSignal chan struct{}) {
	var waitingCount uint8 = 0
	waiting := make(map[string]bool)
	purchase := make(map[string]bool)
	purchaseList := make([]FilePart, 0)
	var lock sync.Mutex
	for _, part := range parts {
		np := FilePart{
			Keys:  make([]string, 0),
			Peers: part.Peers,
		}
		for _, key := range part.Keys {
			this.lock.Lock()
			info := this.files[key]
			if info == nil {
				ninfo := NewFileInfo(key, fUnkown)
				this.files[key] = ninfo
				info = ninfo
			}
			this.lock.Unlock()
			info.lock.Lock()
			if info.state == fPossess {
				info.lock.Unlock()
				continue
			}
			if info.state == fWaiting || info.state == fPurchasing {
				waiting[key] = true
				waitingCount++
			} else {
				waitingCount++
				info.state = fPurchasing
				purchase[key] = true
				np.Keys = append(np.Keys, key)
			}
			info.lock.Unlock()
		}
		//fmt.Println("NP", np.keys)
		if len(np.Keys) > 0 {
			purchaseList = append(purchaseList, np)
		}
	}
	//fmt.Println(purchaseList)
	resultChan := make(chan filePurchaseResult, 5)
	session := NewFilePurchaseSession("RandomID", this, purchaseList, resultChan)
	stopSignal := make(chan struct{}, 1)
	go func() {
		for {
			result, ok := <-resultChan
			if !ok {
				return
			}
			if !result.success {
				// TODO ERR HANDLE
				fmt.Println("TODO")
			} else {
				waitingCount--
				if waitingCount == 0 {
					var ret struct{}
					stopSignal <- ret
					close(stopSignal)

					return
				}
			}
		}
	}()
	session.StartSession()
	for {
		over := false
		select {
		case <-stopSignal:
			over = true
			break
		default:
			time.Sleep(time.Duration(waitRoundCheckDuration))
			lock.Lock()
			//TODO: CHECK WAITINg STATE
			lock.Unlock()
		}
		if over {
			break
		}
	}
	close(resultChan)
	var ret struct{}
	okSignal <- ret
	close(okSignal)
}

func (this *FileSystem) HandleFilePurchaseReq(req peer_list.MessageInfo) {
	var reqobj purchaseReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	//fmt.Println("PURCHASE", reqobj)
	fr := req.From.ToString()
	//TODO: Check Contract
	retfiles := make([][]byte, len(reqobj.Keys))
	for i, key := range reqobj.Keys {
		info := this.files[key]
		info.rwlock.Lock()
		if info.ps[fr] == possess || info.ps[fr] == sending {
			info.rwlock.Unlock()
			retfiles[i] = nil
			continue
		}
		file, err := this.engine.Get(key)
		if err != nil {
			info.rwlock.Unlock()
			retfiles[i] = nil
			continue
		}
		info.ps[fr] = possess
		info.peer = append(info.peer, req.From)
		retfiles[i] = file
		info.rwlock.Unlock()
	}
	//TODO: Invoke Contract
	res := purchaseRes{
		Keys:  reqobj.Keys,
		Files: retfiles,
	}
	resb, _ := json.Marshal(res)
	this.peerList.Reply(req, "HandlePurchaseRes", resb)
}
