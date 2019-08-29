package blockchain

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"vcbb/log"
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
	logger              *log.LoggerInstance
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

func CalculationContractFromAddress(handler *EthBlockChainHandler, addr types.Address, fatherLogger *log.LoggerInstance) (*CalculationContract, error) {
	logger := fatherLogger.GetSubInstance(log.Topic(addr.ToString()))
	logger.Log("Try Create From Address")
	instance, err := NewCalculationProc(common.Address(addr), handler.client)
	if err != nil {
		logger.Err(err.Error())
		return nil, err
	}
	ret := &CalculationContract{
		handler:          handler,
		stopSiganl:       make(chan struct{}),
		Contract:         addr,
		contractInstance: instance,
		logger:           logger,
	}
	logger.Log("Contract Created From Address")
	return ret, nil
}

func NewCalculationContract(
	handler *EthBlockChainHandler,
	contractStateUpdate chan *Answer,
	basicInfo *ContractDeployInfo,
	info *CalculationContractDeployInfo,
	fatherLogger *log.LoggerInstance,
) *CalculationContract {
	fatherLogger.Log("CalcContract Created by: " + handler.account.Id.ToString())
	return &CalculationContract{
		handler:             handler,
		contractStateUpdate: contractStateUpdate,
		basicInfo:           basicInfo,
		info:                info,
		stopSiganl:          make(chan struct{}),
		logger:              fatherLogger, // CAUTION
		//Contract:            types.Address(addr)
		//contractInstance:    instance,,
	}
}

func (this *CalculationContract) Start() (types.Address, error) {
	this.logger.Log("CalculationContract Try Start")
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = this.basicInfo.Value
	auth.GasLimit = this.basicInfo.GasLimit
	this.logger.Log("Try To Get GasPrice")
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		this.logger.Err("Fail To Get GasPrice " + err.Error())
		return types.Address{}, err
	}
	this.logger.Log("GasPrice " + gp.Text(10))
	auth.GasPrice = gp //this.basicInfo.GasPrice
	this.logger.Log("Try To Get Nonce")
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		this.logger.Err("Fail To Get Nonce " + err.Error())
		return types.Address{}, err
	}
	this.logger.Log("Nonce " + strconv.Itoa(int(nonce)))
	auth.Nonce = big.NewInt(int64(nonce))
	this.logger.Log("Try To Deploy Contract")
	addr, _, instance, err := DeployCalculationProc(auth, this.handler.client, this.info.Id, this.info.St, this.info.Fund, this.info.ParticipantCount, this.info.Distribute)
	if err != nil {
		this.logger.Err("Fail To Deploy Contract " + err.Error())
		return types.Address{}, err
	}
	this.logger.Log("Contract Deployed " + addr.String())
	this.logger = this.logger.GetSubInstance(log.Topic(types.Address(addr).ToString()))
	this.Contract = types.Address(addr)
	this.contractInstance = instance
	committedChan := make(chan *CalculationProcCommitted, 5)
	this.logger.Log("Try To Watch Commit")
	this.commitWatcher, err = this.contractInstance.WatchCommitted(&bind.WatchOpts{Context: context.Background()}, committedChan)
	if err != nil {
		this.logger.Err("Watch Commit Fail " + err.Error())
		return types.Address{}, err
	}
	punishedChan := make(chan *CalculationProcPunished, 5)
	this.logger.Log("Try To Watch Punish")
	this.punishWatcher, err = this.contractInstance.WatchPunished(nil /*&bind.WatchOpts{Context: context.Background()}*/, punishedChan)
	if err != nil {
		this.logger.Err("Watch Punish Fail " + err.Error())
		return types.Address{}, err
	}
	terminatedChan := make(chan *CalculationProcTerminated, 5)
	this.logger.Log("Try To Watch Terminate")
	this.terminateWatcher, err = this.contractInstance.WatchTerminated(nil, terminatedChan)
	if err != nil {
		this.logger.Err("Watch Terminate Fail " + err.Error())
		return types.Address{}, err
	}
	go this.watch(committedChan, punishedChan, terminatedChan) //ERROR HANDLE
	this.logger.Log("Strart OK")
	return this.Contract, nil
}

func (this *CalculationContract) watch(committedChan chan *CalculationProcCommitted, punishedChan chan *CalculationProcPunished, terminatedChan chan *CalculationProcTerminated) {
	this.logger.Log("Start Watching")
	for {
		select {
		case commit := <-committedChan:
			this.logger.Log("Receive Commit Answer: " + commit.Ans + " AnsewrHash: " + commit.AnsHash + " Participant: " + types.Address(commit.Participant).ToString())
			var ans [][]string
			err := json.Unmarshal([]byte(commit.Ans), &ans)
			if err != nil {
				this.logger.Err("Commit Msg Unmarshal Fail " + err.Error())
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
			this.logger.Log("Receive Punish " + types.Address(punish.Participant).ToString())
			this.contractStateUpdate <- &Answer{
				Commiter: types.Address(punish.Participant),
			} /*
				case err := <- this.punishWatcher.Err():
				fmt.Println("PUNISH ERR",err)*/
		case terminate := <-terminatedChan:
			this.logger.Log("Receive Terminate Answer: " + terminate.Ans + " Count: " + terminate.Cnt.Text(10))
		case <-this.stopSiganl:
			this.logger.Log("Stop Watching")
			return
		}
	}
}

func (this *CalculationContract) Commit(info *ContractDeployInfo, ans [][]string, ansHash string) error {
	this.logger.Log("Try To Commit")
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	ansb, _ := json.Marshal(ans)
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = info.Value
	auth.GasLimit = info.GasLimit
	this.logger.Log("Try To Get GasPrice")
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		this.logger.Err("Fail To Get GasPrice " + err.Error())
		return err
	}
	this.logger.Log("GasPrice " + gp.Text(10))
	auth.GasPrice = gp //this.basicInfo.GasPrice
	this.logger.Log("Try To Get Nonce")
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		this.logger.Err("Fail To Get Nonce " + err.Error())
		return err
	}
	this.logger.Log("Nonce " + strconv.Itoa(int(nonce)))
	auth.Nonce = big.NewInt(int64(nonce))
	_, err = this.contractInstance.Commit(auth, string(ansb), ansHash)
	if err != nil {
		this.logger.Err("Commit Fail " + err.Error())
		return err
	}
	return nil
}

func (this *CalculationContract) Terminate() error {
	this.logger.Log("Try To Terminate")
	this.handler.lock.Lock()
	defer this.handler.lock.Unlock()
	auth := bind.NewKeyedTransactor(this.handler.account.ECDSAPrivateKey)
	auth.Value = big.NewInt(0)
	auth.GasLimit = this.basicInfo.GasLimit
	this.logger.Log("Try To Get GasPrice")
	gp, err := this.handler.client.SuggestGasPrice(context.Background())
	if err != nil {
		this.logger.Err("Fail To Get GasPrice " + err.Error())
		return err
	}
	this.logger.Log("GasPrice " + gp.Text(10))
	auth.GasPrice = gp //this.basicInfo.GasPrice
	this.logger.Log("Try To Get Nonce")
	nonce, err := this.handler.client.PendingNonceAt(context.Background(), common.Address(this.handler.account.Id))
	if err != nil {
		this.logger.Err("Fail To Get Nonce " + err.Error())
		return err
	}
	this.logger.Log("Nonce " + strconv.Itoa(int(nonce)))
	auth.Nonce = big.NewInt(int64(nonce))
	_, err = this.contractInstance.Terminate(auth)
	if err != nil {
		this.logger.Err("Terminate " + err.Error())
	}
	close(this.contractStateUpdate)
	this.punishWatcher.Unsubscribe()
	this.commitWatcher.Unsubscribe()
	this.terminateWatcher.Unsubscribe()
	this.stopSiganl <- *new(struct{})
	close(this.stopSiganl)
	this.logger.Log("Terminate OK")
	return nil
}
