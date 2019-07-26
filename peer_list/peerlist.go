package peer_list

type PeerList struct {
	instances map[string]*PeerListInstance
}

func (this *PeerList) GetInstance(id string) *PeerListInstance {
	ret := &PeerListInstance{ID: id, PL: this}
	this.instances[id] = ret
	return ret
}

func (this *PeerList) RemoveInstance(id string) {
	delete(this.instances, id)
}
