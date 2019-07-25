package blockchain

import "vcbb/types"

type ComputationContract struct {
	handler BlockChainHandler
}

type ComputationContractUpdate struct {
	NewAnswer map[string][]types.Address
	Punished  map[string][]types.Address
}

func NewComputationContract(handler BlockChainHandler) *ComputationContract {
	return &ComputationContract{handler: handler}
}

func (this *ComputationContract) Start() (types.Address, error) {
	var ret types.Address
	return ret, nil
}

func (this *ComputationContract) Terminate() {

}
