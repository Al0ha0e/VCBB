package peer_list

import (
	"encoding/json"
	"fmt"
	"vcbb/net"

	"github.com/Al0ha0e/vcbb/types"
)

type PeerList struct {
	Address   types.Address
	peers     []types.Address
	instances map[string]*PeerListInstance
	callBack  map[string]func(MessageInfo)
	//channels   map[string]chan MessageInfo
	netService *net.NetSimulator
	bus        chan []byte
	stopSignal chan struct{}
}

func NewPeerList(addr types.Address, ns *net.NetSimulator) *PeerList {
	return &PeerList{
		Address:   addr,
		peers:     make([]types.Address, 0, 10),
		instances: make(map[string]*PeerListInstance),
		callBack:  make(map[string]func(MessageInfo)),
		//channels:   make(map[string]chan MessageInfo),
		netService: ns,
		bus:        make(chan []byte, 10),
		stopSignal: make(chan struct{}, 1),
	}
}

func (this *PeerList) Run() {
	this.netService.RegisterUser(this.Address.ToString(), this.bus)
	go this.Serve()
}

func (this *PeerList) Serve() {
	for {
		select {
		case <-this.stopSignal:
			return
		case msg := <-this.bus:
			var msgobj Message
			json.Unmarshal(msg, &msgobj)
			fmt.Println(msgobj)
			msgobj.Dist--
			if msgobj.Dist > 0 {
				this.BroadCastRPC(msgobj.Method, msgobj.Content, msgobj.Dist)
			}
			msginfo := MessageInfo{
				From:    msgobj.From,
				Content: msgobj.Content,
			}
			if msgobj.Session == "global" {
				method := this.callBack[msgobj.Method]
				if method != nil {
					go method(msginfo)
				}
			} else {
				sess := this.instances[msgobj.Session]
				if sess != nil {
					sess.HandleMsg(msgobj.Method, msginfo)
				}
			}
		}
	}
}

func (this *PeerList) GetInstance(id string) *PeerListInstance {
	ret := NewPeerListInstance(id, this)
	this.instances[id] = ret
	return ret
}

func (this *PeerList) RemoveInstance(id string) {
	delete(this.instances, id)
}

func (this *PeerList) RemoteProcedureCall(to types.Address, method string, msg []byte) error {
	pkg := newMessage(this.Address, to, "global", method, msg, 1)
	pkgb, err := json.Marshal(pkg)
	if err != nil {
		return err
	}
	this.netService.SendMessageTo(to.ToString(), pkgb)
	return nil
}

func (this *PeerList) BroadCastRPC(method string, msg []byte, dist uint8) (*[]types.Address, error) {
	for _, peer := range this.peers {
		pkg := newMessage(this.Address, peer, "global", method, msg, dist)
		pkgb, err := json.Marshal(pkg)
		if err != nil {
			return nil, err
		}
		this.netService.SendMessageTo(peer.ToString(), pkgb)
	}
	return &this.peers, nil
}

func (this *PeerList) AddCallBack(name string, cb func(MessageInfo)) {
	this.callBack[name] = cb
}
