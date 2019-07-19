package blockchain

import (
	"vcbb/types"
)

type BlockChainHandler interface {
	CreateContract(string) (types.Account, error)
	CreateComputationContract() error
}
