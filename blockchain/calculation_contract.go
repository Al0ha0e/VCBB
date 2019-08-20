package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"vcbb/types"

	"github.com/ethereum/go-ethereum/event"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type CalculationContract struct {
	Contract            types.Address
	handler             *EthBlockChainHandler
	contractStateUpdate chan *Answer
	contractInstance    *CalculationProc
	basicInfo           *ContractDeployInfo
	info                *CalculationContractDeployInfo
	commitWatcher       event.Subscription
	punishWatcher       event.Subscription
	stopSiganl          chan struct{}
}

type CalculationContractDeployInfo struct {
	Id               string
	St               *big.Int
	Fund             *big.Int
	ParticipantCount uint8
	Distribute       [8]*big.Int
}

type Answer struct {
	Commiter types.Address
	Ans      [][]string
	AnsHash  string
}

/*
type CalculationContractUpdate struct {
	NewAnswer map[string]*Answer
	Punished  map[string][]types.Address
}*/

func NewCalculationContract(
	handler *EthBlockChainHandler,
	contractStateUpdate chan *Answer,
	basicInfo *ContractDeployInfo,
	info *CalculationContractDeployInfo,
) *CalculationContract {
	return &CalculationContract{
		handler:             handler,
		contractStateUpdate: contractStateUpdate,
		basicInfo:           basicInfo,
		info:                info,
		stopSiganl:          make(chan struct{}),
		//Contract:            types.Address(addr)
		//contractInstance:    instance,,
	}
}

func (this *CalculationContract) Start() (types.Address, error) {
	privateKey, err := crypto.HexToECDSA(this.handler.account.PrivateKey)
	if err != nil {
		return types.Address{}, err
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = this.basicInfo.Value
	auth.GasLimit = this.basicInfo.GasLimit
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		return types.Address{}, err
	}
	auth.GasPrice = gp //this.basicInfo.GasPrice
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		fmt.Println("NONCE", err)
		return types.Address{}, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	addr, _, instance, err := DeployCalculationProc(auth, this.handler.client, this.info.Id, this.info.St, this.info.Fund, this.info.ParticipantCount, this.info.Distribute)
	if err != nil {
		return types.Address{}, err
	}
	this.Contract = types.Address(addr)
	this.contractInstance = instance
	committedChan := make(chan *CalculationProcCommitted)
	this.commitWatcher, err = this.contractInstance.WatchCommitted(&bind.WatchOpts{}, committedChan)
	if err != nil {
		return types.Address{}, err
	}
	punishedChan := make(chan *CalculationProcPunished)
	this.punishWatcher, err = this.contractInstance.WatchPunished(&bind.WatchOpts{}, punishedChan)
	if err != nil {
		return types.Address{}, err
	}
	go this.watch(committedChan, punishedChan) //ERROR HANDLE
	return this.Contract, nil
}

func (this *CalculationContract) watch(committedChan chan *CalculationProcCommitted, punishedChan chan *CalculationProcPunished) {
	for {
		select {
		case commit := <-committedChan:
			var ans [][]string
			err := json.Unmarshal([]byte(commit.Ans), &ans)
			if err != nil {
				continue
			}
			this.contractStateUpdate <- &Answer{
				Commiter: types.Address(commit.Participant),
				Ans:      ans,
				AnsHash:  commit.AnsHash,
			}
		case punish := <-punishedChan:
			this.contractStateUpdate <- &Answer{
				Commiter: types.Address(punish.Participant),
			}
		case <-this.stopSiganl:
			return
		}
	}
}

func (this *CalculationContract) Terminate() {
	close(this.contractStateUpdate)
	this.punishWatcher.Unsubscribe()
	this.commitWatcher.Unsubscribe()
	this.stopSiganl <- *new(struct{})
	close(this.stopSiganl)
}
