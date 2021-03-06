package peer_list

import (
	"vcbb/log"
	"vcbb/types"
)

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
	callBack map[string]func(MessageInfo)
	bus      chan []byte
	logger   *log.LoggerInstance
	//channels map[string]chan MessageInfo
}

func NewPeerListInstance(id string, pl *PeerList, fatherLogger *log.LoggerInstance) *PeerListInstance {
	return &PeerListInstance{
		ID:       id,
		PL:       pl,
		callBack: make(map[string]func(MessageInfo)),
		bus:      make(chan []byte, 10),
		logger:   fatherLogger.GetSubInstance(log.Topic(id)),
		//channels: make(map[string]chan MessageInfo),
	}
}

func (this *PeerListInstance) HandleMsg(meth string, msg MessageInfo) {
	this.logger.Log("Handle Msg From " + msg.From.ToString())
	method := this.callBack[meth]
	if method != nil {
		go method(msg)
	}
}

func (this *PeerListInstance) AddCallBack(name string, cb func(MessageInfo)) {
	this.logger.Log("Add Callback " + name)
	this.callBack[name] = cb
}
func (this *PeerListInstance) GlobalRemoteProcedureCall(to types.Address, method string, msg []byte) error {
	return this.PL.BasicRemoteProcedureCall(to, this.ID, Global, method, msg, 1)
}

func (this *PeerListInstance) Reply(info MessageInfo, metod string, msg []byte) error {
	return this.PL.BasicRemoteProcedureCall(info.From, this.ID, this.ID, metod, msg, 1)
}

func (this *PeerListInstance) SendDataPackTo(to types.Address, pack types.DataPack) {

}

func (this *PeerListInstance) UpdatePunishedPeer(types.Address) {

}

/*
func (this *PeerListInstance) Close() {
	this.PL.RemoveInstance(this.ID)
	for k, v := range this.channels {
		close(v)
		delete(this.channels, k)
	}
}*/
