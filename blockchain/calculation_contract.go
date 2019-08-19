package blockchain

import (
	"vcbb/types"

	"github.com/ethereum/go-ethereum/common"
)

type CalculationContract struct {
	Contract            types.Address
	handler             EthBlockChainHandler
	contractStateUpdate chan *CalculationContractUpdate
	contractInstance    *CalculationProc
}

type Answer struct {
	Commiters []types.Address
	Ans       [][]string
	AnsHash   string
}

type CalculationContractUpdate struct {
	NewAnswer map[string]*Answer
	Punished  map[string][]types.Address
}

func NewCalculationContract(handler EthBlockChainHandler, contractStateUpdate chan *CalculationContractUpdate) (*CalculationContract, error) {
	var address types.Address
	instance, err := NewCalculationProc(common.Address(address), handler.client)
	if err != nil {
		return nil, err
	}
	return &CalculationContract{
		Contract:            address,
		handler:             handler,
		contractStateUpdate: contractStateUpdate,
		contractInstance:    instance,
	}, nil
}

func (this *CalculationContract) Start() (types.Address, error) {
	var ret types.Address
	return ret, nil
}

func (this *CalculationContract) Terminate() {
	close(this.contractStateUpdate)
}
