package peer_list

import (
	"encoding/json"
	"strconv"
	"vcbb/net"

	"vcbb/log"
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
	netService net.NetService
	bus        chan []byte
	stopSignal chan struct{}
	logger     *log.LoggerInstance
}

func NewPeerList(addr types.Address, ns net.NetService, logSystem *log.LogSystem) *PeerList {
	return &PeerList{
		Address:   addr,
		Peers:     make([]types.Address, 0, 10),
		instances: make(map[string]*PeerListInstance),
		callBack:  make(map[string]func(MessageInfo)),
		//channels:   make(map[string]chan MessageInfo),
		netService: ns,
		bus:        make(chan []byte, 10),
		stopSignal: make(chan struct{}, 1),
		logger:     logSystem.GetInstance(log.Topic("PeerList")),
	}
}

func (this *PeerList) Run() {
	this.logger.Log("PeerList Start")
	this.netService.RegisterUser(this.Address.ToString(), this.bus)
	go this.Serve()
}

func (this *PeerList) Serve() {
	this.logger.Log("Serving")
	for {
		select {
		case <-this.stopSignal:
			this.logger.Log("Stop")
			return
		case msg := <-this.bus:
			this.logger.Log("Receive Message")
			var msgobj Message
			err := json.Unmarshal(msg, &msgobj)
			if err != nil {
				this.logger.Log("Fail To Unmarshal Message")
				continue
			}
			this.logger.Log(
				"Message From: " + msgobj.From.ToString() +
					" From Session: " + msgobj.FromSession +
					" To Session " + msgobj.ToSession +
					" To Method: " + msgobj.Method +
					" Distance " + strconv.Itoa(int(msgobj.Dist)))
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
	this.logger.Log("Get PeerList Instance " + id)
	ret := NewPeerListInstance(id, this, this.logger)
	this.instances[id] = ret
	return ret
}

func (this *PeerList) RemoveInstance(id string) {
	this.logger.Log("Remove Instance " + id)
	delete(this.instances, id)
}

func (this *PeerList) RemoteProcedureCall(to types.Address, toSession, method string, msg []byte) error {
	return this.BasicRemoteProcedureCall(to, Global, toSession, method, msg, 1)
}

func (this *PeerList) BroadCastRPC(method string, msg []byte, dist uint8) ([]types.Address, error) {
	this.logger.Log("BroadCast " + method)
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
	this.logger.Log("Add Callback " + name)
	this.callBack[name] = cb
}

func (this *PeerList) BasicRemoteProcedureCall(to types.Address, frsess, tosess, method string, msg []byte, dist uint8) error {
	this.logger.Log("Try To Send Mseeage To: " + to.ToString() +
		" To Session: " + tosess +
		" To Method: " + method +
		" Distance: " + strconv.Itoa(int(dist)))
	pkg := newMessage(this.Address, to, frsess, tosess, method, msg, dist)
	pkgb, err := json.Marshal(pkg)
	if err != nil {
		return err
	}
	this.netService.SendMessageTo(to.ToString(), pkgb)
	return nil
}
