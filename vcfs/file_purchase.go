package vcfs

import (
	"encoding/json"

	"github.com/Al0ha0e/vcbb/peer_list"
)

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
	this.peerList.Reply(req, "", resb)
}
