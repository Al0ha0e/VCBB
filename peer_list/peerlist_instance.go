package peer_list

import "vcbb/types"

type PeerListInstance struct {
	ID          string
	PL          *PeerList
	DataReq     chan MessageInfo
	DataRecv    chan MessageInfo
	MetaDataReq chan MessageInfo
}

func (this *PeerListInstance) Init(datareq chan MessageInfo, datarecv chan MessageInfo, metadatareq chan MessageInfo) {
	this.DataReq = datareq
	this.DataRecv = datarecv
	this.MetaDataReq = metadatareq
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
	if this.DataReq != nil {
		close(this.DataReq)
	}
	if this.DataRecv != nil {
		close(this.DataRecv)
	}
	if this.MetaDataReq != nil {
		close(this.MetaDataReq)
	}
}
