package peer_list

import "vcbb/types"

type PeerList struct {
	Address   types.Address
	instances map[string]*PeerListInstance
	channels  map[string]chan MessageInfo
}

func (this *PeerList) GetInstance(id string) *PeerListInstance {
	ret := &PeerListInstance{ID: id, PL: this}
	this.instances[id] = ret
	return ret
}

func (this *PeerList) RemoveInstance(id string) {
	delete(this.instances, id)
}

func (this *PeerList) RemoteProcedureCall(to types.Address, method string, msg []byte) error {
	return nil
}

func (this *PeerList) BroadCastRPC(method string, msg []byte) error {
	return nil
}

func (this *PeerList) AddChannel(name string, ch chan MessageInfo) {
	this.channels[name] = ch
}
