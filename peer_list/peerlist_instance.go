package peer_list

import "vcbb/types"

const (
	DataReq            = "DataReq"     //DATA STORE SEND TO DATA PROVIDER TO GET DATA
	DataRecv           = "DataRecv"    //MSG SEND TO PROVIDER BY PEERLIST TO INFORM THE END OF A FILE TRANSPORT
	MetaDataReq        = "MetaDataReq" //SLAVE SEND IT TO MASTER TO GET METADATA
	MetaDataRes        = "MetaDataRes"
	InfoReq            = "TrackReq" //SEND TO TRACKER TO GET DATA POSITION
	InfoRes            = "InfoRes"
	SeekReceiverReq    = "SeekReceiverReq" //DATA PROVIDER SEND TO SEEK FOR RECEIVER
	SeekParticipantReq = "SeekParticipantReq"
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
func (this *PeerListInstance) RemoteProcedureCall(to types.Address, session, method string, msg []byte) error {
	return nil
}
func (this *PeerListInstance) BroadCastRPC(session, method string, msg []byte) error {
	return nil
}

/*
func (this *PeerListInstance) BroadCast([]byte) error {
	return nil
}

func (this *PeerListInstance) SendMsgTo(to types.Address, msg []byte) {

}
*/
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
