package blockchain

import (
	"math/big"
	"vcbb/types"

	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockChainHandler interface {
	CreateContract()
}

type EthBlockChainHandler struct {
	client  *ethclient.Client
	account types.Account
}

func NewEthBlockChainHandler(url string, account types.Account) (*EthBlockChainHandler, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthBlockChainHandler{
		client:  cli,
		account: account,
	}, nil
}

type ContractDeployInfo struct {
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
}
