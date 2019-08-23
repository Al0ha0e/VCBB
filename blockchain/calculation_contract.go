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
	terminateWatcher    event.Subscription
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

func CalculationContractFromAddress(handler *EthBlockChainHandler, addr types.Address) (*CalculationContract, error) {
	instance, err := NewCalculationProc(common.Address(addr), handler.client)
	if err != nil {
		//handler.logger.Log("CalcContractFromAddress " + addr.ToString() + "Err: " + err.Error())
		return nil, err
	}
	//handler.logger.Log("CalcContractFromAddress " + addr.ToString() + "Created")
	return &CalculationContract{
		handler:          handler,
		stopSiganl:       make(chan struct{}),
		Contract:         addr,
		contractInstance: instance,
	}, nil
}

func NewCalculationContract(
	handler *EthBlockChainHandler,
	contractStateUpdate chan *Answer,
	basicInfo *ContractDeployInfo,
	info *CalculationContractDeployInfo,
) *CalculationContract {
	//handler.logger.Log("CalcContract Created by: " + handler.account.Id.ToString())
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
	//this.handler.logger.Log("CalculationContract Try Start")
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = this.basicInfo.Value
	auth.GasLimit = this.basicInfo.GasLimit
	fmt.Println("TRY GASPRICE")
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		return types.Address{}, err
	}
	fmt.Println("GP OK", gp)
	auth.GasPrice = gp //this.basicInfo.GasPrice
	fmt.Println("TRY NONCE")
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		fmt.Println("NONCE", err)
		return types.Address{}, err
	}
	fmt.Println("NONCE OK", nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	addr, _, instance, err := DeployCalculationProc(auth, this.handler.client, this.info.Id, this.info.St, this.info.Fund, this.info.ParticipantCount, this.info.Distribute)
	if err != nil {
		return types.Address{}, err
	}
	fmt.Println("OK DEPLOY", addr)
	this.Contract = types.Address(addr)
	this.contractInstance = instance
	committedChan := make(chan *CalculationProcCommitted, 5)
	this.commitWatcher, err = this.contractInstance.WatchCommitted(&bind.WatchOpts{Context: context.Background()}, committedChan)
	if err != nil {
		fmt.Println("COMMITERR", err)
		return types.Address{}, err
	}
	punishedChan := make(chan *CalculationProcPunished, 5)
	this.punishWatcher, err = this.contractInstance.WatchPunished(nil /*&bind.WatchOpts{Context: context.Background()}*/, punishedChan)
	if err != nil {
		fmt.Println("PUNISH ERR", err)
		return types.Address{}, err
	}
	terminatedChan := make(chan *CalculationProcTerminated, 5)
	this.terminateWatcher, err = this.contractInstance.WatchTerminated(nil, terminatedChan)
	if err != nil {
		fmt.Println("TERMINATE ERR", err)
		return types.Address{}, err
	}
	go this.watch(committedChan, punishedChan, terminatedChan) //ERROR HANDLE
	return this.Contract, nil
}

func (this *CalculationContract) watch(committedChan chan *CalculationProcCommitted, punishedChan chan *CalculationProcPunished, terminatedChan chan *CalculationProcTerminated) {
	fmt.Println("WATCHING")
	for {
		select {
		case commit := <-committedChan:
			fmt.Println("RECEIVE COMMIT", commit.Ans, commit.AnsHash, commit.Participant)
			var ans [][]string
			err := json.Unmarshal([]byte(commit.Ans), &ans)
			if err != nil {
				continue
			}
			this.contractStateUpdate <- &Answer{
				Commiter: types.Address(commit.Participant),
				Ans:      ans,
				AnsHash:  commit.AnsHash,
			} /*
				case err := <-this.commitWatcher.Err():
					fmt.Println("COMMIT ERR", err)*/
		case punish := <-punishedChan:
			fmt.Println("RECEIVE PUNISH", punish)
			this.contractStateUpdate <- &Answer{
				Commiter: types.Address(punish.Participant),
			} /*
				case err := <- this.punishWatcher.Err():
				fmt.Println("PUNISH ERR",err)*/
		case terminate := <-terminatedChan:
			fmt.Println("TERMINATE", terminate)
		case <-this.stopSiganl:
			fmt.Println("CONTRACT STOP")
			return
		}
	}
}

func (this *CalculationContract) Commit(info *ContractDeployInfo, ans [][]string, ansHash string) error {
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	ansb, _ := json.Marshal(ans)
	fmt.Println(this.handler.account)
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = info.Value
	auth.GasLimit = info.GasLimit
	//fmt.Println("TRY GASPRICE")
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("GAS PRICE ERR", err)
		return err
	}
	//fmt.Println("GP OK", gp)
	auth.GasPrice = gp //this.basicInfo.GasPrice
	//fmt.Println("TRY NONCE")
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		fmt.Println("NONCE", err)
		return err
	}
	//fmt.Println("NONCE OK", nonce)
	auth.Nonce = big.NewInt(int64(nonce))
	fmt.Println("AUTH", auth)
	_, err = this.contractInstance.Commit(auth, string(ansb), ansHash)
	return err
}

func (this *CalculationContract) Terminate() error {
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = big.NewInt(0)
	auth.GasLimit = this.basicInfo.GasLimit
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	auth.GasPrice = gp //this.basicInfo.GasPrice
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		fmt.Println("NONCE", err)
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	this.contractInstance.Terminate(auth)
	close(this.contractStateUpdate)
	this.punishWatcher.Unsubscribe()
	this.commitWatcher.Unsubscribe()
	this.terminateWatcher.Unsubscribe()
	this.stopSiganl <- *new(struct{})
	close(this.stopSiganl)
	return nil
}
