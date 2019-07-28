package blockchain

import "github.com/Al0ha0e/vcbb/types"

type ComputationContract struct {
	handler             BlockChainHandler
	contractStateUpdate chan *ComputationContractUpdate
}

type ComputationContractUpdate struct {
	NewAnswer map[string][]types.Address
	Punished  map[string][]types.Address
}

func NewComputationContract(handler BlockChainHandler, contractStateUpdate chan *ComputationContractUpdate) *ComputationContract {
	return &ComputationContract{
		handler:             handler,
		contractStateUpdate: contractStateUpdate,
	}
}

func (this *ComputationContract) Start() (types.Address, error) {
	var ret types.Address
	return ret, nil
}

func (this *ComputationContract) Terminate() {
	close(this.contractStateUpdate)
}
