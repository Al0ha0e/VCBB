package blockchain

import (
	"math/big"
	"sync"
	"vcbb/log"
	"vcbb/types"

	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockChainHandler interface {
	CreateContract()
}

type EthBlockChainHandler struct {
	blockchainUrl string
	client        *ethclient.Client
	account       *types.Account
	lock          sync.Mutex
	logger        *log.LoggerInstance
}

func NewEthBlockChainHandler(url string, account *types.Account /*, logSystem *log.LogSystem*/) (*EthBlockChainHandler, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthBlockChainHandler{
		blockchainUrl: url,
		account:       account,
		client:        cli,
		//logger:        logSystem.GetInstance(log.Topic("EthBlockchainHandler")),
	}, nil
}

type ContractDeployInfo struct {
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
}
