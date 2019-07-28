package blockchain

import "github.com/Al0ha0e/vcbb/types"

type DataTransportationContract struct {
	handler BlockChainHandler
}

func NewDataTransportationContract(handler BlockChainHandler) *DataTransportationContract {
	return &DataTransportationContract{handler: handler}
}

func (this *DataTransportationContract) Start() (types.Address, error) {
	var ret types.Address
	return ret, nil
}

func (this *DataTransportationContract) Terminate() {

}
