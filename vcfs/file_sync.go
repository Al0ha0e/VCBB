package vcfs

import (
	"encoding/json"
	"fmt"

	"github.com/Al0ha0e/vcbb/peer_list"
	"github.com/Al0ha0e/vcbb/types"
)

func (this *FileSystem) SyncFile(keys []string, siz []uint64, peers []types.Address) error {
	for _, key := range keys {
		this.rwlock.RLock()
		info := this.files[key]
		if info == nil {
			this.rwlock.RUnlock()
			return fmt.Errorf("file %s not exist", key)
		}
		this.rwlock.RUnlock()
		info.lock.Lock()
		for _, rpeer := range peers {
			peer := rpeer.ToString()
			sta := info.ps[peer]
			if sta == unkown {
				info.ps[peer] = waiting
			}
		}
		info.lock.Unlock()
	}
	req := syncReq{
		Keys: keys,
		Size: siz,
	}
	reqb, _ := json.Marshal(req)
	for _, peer := range peers {
		this.peerList.RemoteProcedureCall(peer, peer_list.Global, "HandleSyncReq", reqb)
	}
	return nil
}

func (this *FileSystem) HandleSyncReq(req peer_list.MessageInfo) {
	var reqobj syncReq
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	fr := req.From.ToString()
	retkeys := make([]string, 0)
	for i, id := range reqobj.Keys {
		this.lock.Lock()
		info := this.files[id]
		if info == nil {
			ninfo := NewFileInfo(id, fWaiting)
			ninfo.peer = append(ninfo.peer, req.From) //[0] = req.From
			ninfo.ps[fr] = possess
			this.files[id] = ninfo
			info = ninfo
		}
		this.lock.Unlock()
		info.lock.Lock()
		prstate := info.ps[fr]
		if prstate != possess {
			info.ps[fr] = possess
			info.peer = append(info.peer, req.From)
		}
		if !(info.state == fPossess) {
			if !this.engine.CanSet(reqobj.Size[i]) {
				info.lock.Unlock()
				continue
			}
			info.ps[fr] = sending
			info.state = fWaiting
			retkeys = append(retkeys, id)
		}
		info.lock.Unlock()
	}
	res := syncRes{
		Keys: retkeys,
	}
	resb, err := json.Marshal(res)
	if err != nil {
		return
	}
	this.peerList.Reply(req, "HandleSyncRes", resb)
}

func (this *FileSystem) HandleSyncRes(req peer_list.MessageInfo) {
	var reqobj syncRes
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	retfiles := make([][]byte, len(reqobj.Keys))
	fr := req.From.ToString()
	for i, key := range reqobj.Keys {
		info := this.files[key]
		info.lock.Lock()
		if info.ps[fr] != waiting {
			info.lock.Unlock()
			retfiles[i] = nil
			continue
		}
		file, err := this.engine.Get(key)
		//fmt.Println("FL", file)
		if err != nil {
			info.lock.Unlock()
			retfiles[i] = nil
			continue
		}
		info.ps[fr] = possess
		info.peer = append(info.peer, req.From)
		retfiles[i] = file
		info.lock.Unlock()
	}
	res := filePack{
		Keys:  reqobj.Keys,
		Files: retfiles,
	}
	//fmt.Println("RES", res)
	resb, _ := json.Marshal(res)
	this.peerList.Reply(req, "HandleSyncFileArrive", resb)
}

func (this *FileSystem) HandleSyncFileArrive(req peer_list.MessageInfo) {
	var reqobj filePack
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	//fmt.Println(reqobj)
	fr := req.From.ToString()
	for i, key := range reqobj.Keys {
		info := this.files[key]
		if info != nil {
			info.lock.Lock()
			if info.ps[fr] != sending {
				info.lock.Unlock()
				continue
			}
			info.ps[fr] = possess
			file := reqobj.Files[i]
			//TODO: FILE CHECK
			if file != nil {
				err := this.engine.Set(key, file)
				if err != nil {
					info.lock.Unlock()
					continue
				}
				info.state = fPossess
				// TODO: DISPATCH FILE ARRIVE EVENT
			}
			info.lock.Unlock()
		}
	}
}
