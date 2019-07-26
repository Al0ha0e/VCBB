package data

/*TODO:
* change participant state
*
 */

import (
	"encoding/json"
	"fmt"
	"vcbb/peer_list"
	"vcbb/types"
)

type DataTransportState uint8

const (
	Preparing DataTransportState = iota
	Transporting
	Finished
)

type PtState uint8

const (
	GotMeta PtState = iota
	Received
)

type DataTransportSession struct {
	ID                  string
	State               DataTransportState
	PeerList            *peer_list.PeerListInstance
	DataSource          types.DataSource
	ParticipantState    map[types.Address]PtState
	DataHash            []string
	DataDistribute      []uint64
	DataReceivers       map[string][]types.Address
	DataReq             chan peer_list.MessageInfo
	DataRecv            chan peer_list.MessageInfo
	DataTransportResult chan *DataTransportMeta
	TerminateSignal     chan struct{}
}

type DataTransportMeta struct {
	DataSource    types.DataSource
	DataHash      []string
	DataReceivers map[string][]types.Address
}

func NewDataTransportSession(id string, peerlist *peer_list.PeerList, datasource types.DataSource) *DataTransportSession {
	return &DataTransportSession{
		ID:         id,
		State:      Preparing,
		PeerList:   peerlist.GetInstance(id),
		DataSource: datasource,
		DataHash:   datasource.GetHashList(),
		//DataDistribute: make([]uint64, 0),
		//DataReceivers:  make(map[string][]types.Address),
		//DataReq:             make(chan peer_list.MessageInfo, 1),
		//DataRecv:            make(chan peer_list.MessageInfo, 1),
		//DataTransportResult: make(chan *DataTransportMeta, 1),
		//TerminateSignal:     make(chan struct{}, 2),
	}
}

func (this *DataTransportSession) Init() {
	this.ParticipantState = make(map[types.Address]PtState)
	this.DataDistribute = make([]uint64, 0)
	this.DataReceivers = make(map[string][]types.Address)
	this.DataReq = make(chan peer_list.MessageInfo, 1)
	this.DataRecv = make(chan peer_list.MessageInfo, 1)
	this.DataTransportResult = make(chan *DataTransportMeta, 1)
	this.TerminateSignal = make(chan struct{}, 2)
	this.PeerList.Init(this.DataReq, this.DataRecv, nil)
}

func (this *DataTransportSession) StartSession() (chan *DataTransportMeta, error) {
	this.Init()
	req, err := this.getDataTransportReq()
	if err != nil {
		return nil, err
	}
	this.PeerList.BroadCast(req)
	this.State = Transporting
	go this.handleDataReq()
	go this.updateDataReceiver()
	this.DataTransportResult = make(chan *DataTransportMeta, 1)
	return this.DataTransportResult, nil
}

func (this *DataTransportSession) getDataTransportReq() ([]byte, error) {
	var req dataTransportReq
	req.Metadata = this.DataHash
	reqb, err := json.Marshal(req)
	return reqb, err
}

func (this *DataTransportSession) handleDataReq() {
	var l, r uint64
	for {
		select {
		case req := <-this.DataReq:
			if this.ParticipantState[req.From] != GotMeta {
				continue
			}
			var reqobj dataTransportRes
			err := json.Unmarshal(req.Content, &reqobj)
			if err != nil {
				continue
			}
			var dataPack types.DataPack
			dataPack, l, r, err = this.getDataByAmount(reqobj.Amount, l, r)
			this.PeerList.SendDataPackTo(req.From, dataPack)
		case <-this.TerminateSignal:
			return
		}
	}
}

func (this *DataTransportSession) getDataByAmount(amount, l, r uint64) (types.DataPack, uint64, uint64, error) {
	tl, tr := l, r
	if amount == 0 {
		return nil, l, r, fmt.Errorf("amount of data must be positive")
	}
	pl := uint64(len(this.DataHash))
	if amount > pl {
		amount = pl
	}
	var ret types.DataPack
	for i := l - 1; i >= 0 && amount > 0; i++ {
		tl--
		amount--
		this.DataDistribute[i]++
		single, err := this.DataSource.GetSingle(i)
		if err != nil {
			return nil, l, r, err
		}
		ret = append(ret, single)
	}
	for i := r; i < pl && amount > 0; i++ {
		tr++
		amount--
		this.DataDistribute[i]++
		single, err := this.DataSource.GetSingle(i)
		if err != nil {
			return nil, l, r, err
		}
		ret = append(ret, single)
	}
	if amount > 0 {
		tl = l
		tr = l + 1
		for i := l; i < r && amount > 0; i++ {
			tr++
			amount--
			this.DataDistribute[i]++
			single, err := this.DataSource.GetSingle(i)
			if err != nil {
				return nil, l, r, err
			}
			ret = append(ret, single)
		}
	}
	return ret, tl, tr, nil
}

func (this *DataTransportSession) updateDataReceiver() {
	for {
		select {
		case res := <-this.DataRecv:
			var resobj dataReceivedRes
			err := json.Unmarshal(res.Content, &resobj)
			if err != nil {
				continue
			}
			for _, i := range resobj.DataList {
				this.DataReceivers[this.DataHash[i]] = append(this.DataReceivers[this.DataHash[i]], res.From)
			}
		case <-this.TerminateSignal:
			return
		}
	}
}

func (this *DataTransportSession) Terminate() {
	this.State = Finished
	this.PeerList.Close()
	retmeta := &DataTransportMeta{
		DataSource:    this.DataSource,
		DataHash:      this.DataHash,
		DataReceivers: this.DataReceivers}
	this.DataTransportResult <- retmeta
	close(this.DataTransportResult)
	this.TerminateSignal <- *new(struct{})
	this.TerminateSignal <- *new(struct{})
	close(this.TerminateSignal)
}
