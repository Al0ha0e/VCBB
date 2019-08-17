package peer_list

import (
	"encoding/json"
	"vcbb/net"

	"vcbb/types"
)

const (
	Global = "global"
)

type PeerList struct {
	Address   types.Address
	Peers     []types.Address
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
		Peers:     make([]types.Address, 0, 10),
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
			//fmt.Println(msgobj)
			msgobj.Dist--
			if msgobj.Dist > 0 {
				this.BroadCastRPC(msgobj.Method, msgobj.Content, msgobj.Dist)
			}
			msginfo := NewMessageInfo(msgobj.From, msgobj.FromSession, msgobj.Content)
			if msgobj.ToSession == Global {
				method := this.callBack[msgobj.Method]
				if method != nil {
					go method(msginfo)
				}
			} else {
				sess := this.instances[msgobj.ToSession]
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

func (this *PeerList) RemoteProcedureCall(to types.Address, toSession, method string, msg []byte) error {
	return this.BasicRemoteProcedureCall(to, Global, toSession, method, msg, 1)
}

func (this *PeerList) BroadCastRPC(method string, msg []byte, dist uint8) ([]types.Address, error) {
	for _, peer := range this.Peers {
		err := this.BasicRemoteProcedureCall(peer, Global, Global, method, msg, dist)
		if err != nil {
			return nil, err
		}
	}
	return this.Peers, nil
}

func (this *PeerList) Reply(info MessageInfo, method string, msg []byte) error {
	return this.RemoteProcedureCall(info.From, info.FromSession, method, msg)
}

func (this *PeerList) AddCallBack(name string, cb func(MessageInfo)) {
	this.callBack[name] = cb
}

func (this *PeerList) BasicRemoteProcedureCall(to types.Address, frsess, tosess, method string, msg []byte, dist uint8) error {
	pkg := newMessage(this.Address, to, frsess, tosess, method, msg, dist)
	pkgb, err := json.Marshal(pkg)
	if err != nil {
		return err
	}
	this.netService.SendMessageTo(to.ToString(), pkgb)
	return nil
}
