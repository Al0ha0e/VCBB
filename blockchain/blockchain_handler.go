package blockchain

import (
	"vcbb/types"

	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockChainHandler interface {
	CreateContract()
}

type EthBlockChainHandler struct {
	client  *ethclient.Client
	address types.Address
}

func NewEthBlockChainHandler(url string, address types.Address) (*EthBlockChainHandler, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthBlockChainHandler{
		client:  cli,
		address: address,
	}, nil
}
