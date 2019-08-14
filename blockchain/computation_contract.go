package blockchain

import "vcbb/types"

type ComputationContract struct {
	Contract            types.Address
	handler             BlockChainHandler
	contractStateUpdate chan *ComputationContractUpdate
}

type Answer struct {
	Commiters []types.Address
	Ans       [][]string
}

type ComputationContractUpdate struct {
	NewAnswer map[string]*Answer
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
