package vcfs

import (
	"encoding/json"
	"fmt"

	"github.com/Al0ha0e/vcbb/peer_list"
)

func (this *FileSystem) FetchFiles(parts []filePart, okSignal chan struct{}) {
	var waitingCount uint8 = 0
	var waiting map[string]bool
	var purchase map[string]bool
	purchaseList := make([]filePart, 0)
	for _, part := range parts {
		np := filePart{
			keys:  make([]string, 0),
			peers: part.peers,
		}
		for _, key := range part.keys {
			info := this.files[key]
			info.rwlock.RLock()
			if info.state == fPossess {
				info.rwlock.RUnlock()
				continue
			}
			if info.state == fWaiting || info.state == fPurchasing {
				waiting[key] = true
				waitingCount++
			} else {
				waitingCount++
				info.state = fPurchasing
				purchase[key] = true
				np.keys = append(np.keys, key)
			}
			info.rwlock.RUnlock()
		}
		if len(np.keys) > 0 {
			purchaseList = append(purchaseList, np)
		}
	}
	resultChan := make(chan filePurchaseResult, 5)
	session := NewFilePurchaseSession("", this, purchaseList, resultChan)
	go func() {
		for {
			result, ok := <-resultChan
			if !ok {
				return
			}
			if !result.success {
				fmt.Println("TODO")
			} else {
				waitingCount--
				if waitingCount == 0 {

				}
			}
		}
	}()
	go func() {
		///TODO ROUND CHECK
	}()
	session.StartSession()
}

func (this *FileSystem) HandleFilePurchaseReq(req peer_list.MessageInfo) {
	var reqobj purchaseReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	fr := req.From.ToString()
	//TODO: Check Contract
	retfiles := make([][]byte, len(reqobj.Keys))
	for _, key := range reqobj.Keys {
		info := this.files[key]
		info.rwlock.Lock()
		if info.ps[fr] == possess || info.ps[fr] == sending {
			info.rwlock.Unlock()
			retfiles = append(retfiles, nil)
			continue
		}
		file, err := this.engine.Get(key)
		if err != nil {
			info.rwlock.Unlock()
			retfiles = append(retfiles, nil)
			continue
		}
		info.ps[fr] = possess
		info.peer = append(info.peer, req.From)
		retfiles = append(retfiles, file)
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
