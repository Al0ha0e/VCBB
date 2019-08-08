package vcfs

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
)

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
			ninfo := NewFileInfo(id, false)
			ninfo.peer[0] = req.From
			ninfo.peerState[fr] = possess
			this.files[id] = ninfo
			info = ninfo
		}
		this.lock.Unlock()
		info.lock.Lock()
		prstate := info.peerState[fr]
		if prstate != possess {
			info.peerState[fr] = possess
			info.peer = append(info.peer, req.From)
		}
		if !info.local {
			if !this.engine.CanSet(reqobj.Size[i]) {
				info.lock.Unlock()
				continue
			}
			info.peerState[fr] = sending
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
	for _, key := range reqobj.Keys {
		info := this.files[key]
		info.lock.Lock()
		if info.peerState[fr] != waiting {
			info.lock.Unlock()
			retfiles = append(retfiles, nil)
			continue
		}
		file, err := this.engine.Get(key)
		if err != nil {
			info.lock.Unlock()
			retfiles = append(retfiles, nil)
			continue
		}
		info.peerState[fr] = possess
		info.peer = append(info.peer, req.From)
		retfiles = append(retfiles, file)
		info.lock.Unlock()
	}
	res := filePack{
		Keys:  reqobj.Keys,
		Files: retfiles,
	}
	resb, _ := json.Marshal(res)
	this.peerList.Reply(req, "HandleSyncFileArrive", resb)
}

func (this *FileSystem) HandleSyncFileArrive(req peer_list.MessageInfo) {
	var reqobj filePack
	err := json.Unmarshal(req.Content, &reqobj)
	if err != nil {
		return
	}
	fr := req.From.ToString()
	for i, key := range reqobj.Keys {
		info := this.files[key]
		if info != nil {
			info.lock.Lock()
			if info.peerState[fr] != sending {
				info.lock.Unlock()
				continue
			}
			info.peerState[fr] = possess
			file := reqobj.Files[i]
			if file != nil {
				err := this.engine.Set(key, file)
				if err != nil {
					info.lock.Unlock()
					continue
				}
				info.local = true
			}
			info.lock.Unlock()
		}
	}
}
