package peer_list

import "vcbb/types"

const (
	DataReq     = "DataReq"
	DataRecv    = "DataRecv"
	MetaDataReq = "MetaDataReq"
	InfoReq     = "TrackReq"
)

type PeerListInstance struct {
	ID       string
	PL       *PeerList
	channels map[string]chan MessageInfo
	//DataReq     chan MessageInfo
	//DataRecv    chan MessageInfo
	//MetaDataReq chan MessageInfo
}

/*
func (this *PeerListInstance) Init(datareq chan MessageInfo, datarecv chan MessageInfo, metadatareq chan MessageInfo) {
	this.DataReq = datareq
	this.DataRecv = datarecv
	this.MetaDataReq = metadatareq
}*/

func (this *PeerListInstance) AddChannel(name string, ch chan MessageInfo) {
	this.channels[name] = ch
}

func (this *PeerListInstance) BroadCast([]byte) error {
	return nil
}

func (this *PeerListInstance) SendMsgTo(to types.Address, msg []byte) {

}

func (this *PeerListInstance) SendDataPackTo(to types.Address, pack types.DataPack) {

}

func (this *PeerListInstance) UpdatePunishedPeers(map[string][]types.Address) {

}

func (this *PeerListInstance) Close() {
	this.PL.RemoveInstance(this.ID)
	for k, v := range this.channels {
		close(v)
		delete(this.channels, k)
	}
}
