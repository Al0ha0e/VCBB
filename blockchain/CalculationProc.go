// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CalculationProcABI is the input ABI used to generate the binding from.
const CalculationProcABI = "[{\"constant\":false,\"inputs\":[],\"name\":\"terminate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"answerHash\",\"type\":\"string\"}],\"name\":\"commit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"st\",\"type\":\"uint256\"},{\"name\":\"fund\",\"type\":\"uint256\"},{\"name\":\"participantCount\",\"type\":\"uint8\"},{\"name\":\"distribute\",\"type\":\"uint256[8]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"ansHash\",\"type\":\"string\"}],\"name\":\"committed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"punished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"ans\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"cnt\",\"type\":\"uint256\"}],\"name\":\"terminated\",\"type\":\"event\"}]"

// CalculationProcBin is the compiled bytecode used for deploying new contracts.
const CalculationProcBin = `6080604052604051620015c1380380620015c183398101806040526101808110156200002a57600080fd5b8101908080516401000000008111156200004357600080fd5b828101905060208101848111156200005a57600080fd5b81518560018202830111640100000000821117156200007857600080fd5b50509291906020018051906020019092919080519060200190929190805190602001909291909190505060008260ff16118015620000bd5750600860ff168260ff1611155b151562000132576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f696e76616c69642072657761726465645061727469636970616e74436f756e7481525060200191505060405180910390fd5b60648310151515620001ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6d61737465722066756e64206e6f7420656e6f7567680000000000000000000081525060200191505060405180910390fd5b600083905060008090505b8360ff168160ff161015620001ef57828160ff16600881101515620001d857fe5b6020020151820191508080600101915050620001b7565b50348114151562000268576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600c8152602001807f696e76616c69642066756e64000000000000000000000000000000000000000081525060200191505060405180910390fd5b85600c90805190602001906200028092919062000386565b508460038190555082600260016101000a81548160ff021916908360ff160217905550816004906008620002b69291906200040d565b506001600260006101000a81548160ff021916908360ff16021790555033600d60006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555083600f6000600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505050505050506200047a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620003c957805160ff1916838001178555620003fa565b82800160010185558215620003fa579182015b82811115620003f9578251825591602001919060010190620003dc565b5b50905062000409919062000452565b5090565b82600881019282156200043f579160200282015b828111156200043e57825182559160200191906001019062000421565b5b5090506200044e919062000452565b5090565b6200047791905b808211156200047357600081600090555060010162000459565b5090565b90565b611137806200048a6000396000f3fe608060405260043610610046576000357c0100000000000000000000000000000000000000000000000000000000900480630c08bf88146100485780639867db741461005f575b005b34801561005457600080fd5b5061005d61011a565b005b6101186004803603602081101561007557600080fd5b810190808035906020019064010000000081111561009257600080fd5b8201836020820111156100a457600080fd5b803590602001918460018302840111640100000000831117156100c657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610ae4565b005b6003544210151515610194576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f707265706172696e67000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156101f9576101f433610f50565b610ae2565b6000809050606060008090505b60015481101561050f576000808281548110151561022057fe5b906000526020600020906002020160010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600e60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156102ac5750610502565b60106000838154811015156102bd57fe5b906000526020600020906002020160000160405180828054600181600116156101000203166002900480156103295780601f10610307576101008083540402835291820191610329565b820191906000526020600020905b815481529060010190602001808311610315575b505091505090815260200160405180910390208190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050600060106000848154811015156103b357fe5b9060005260206000209060020201600001604051808280546001816001161561010002031660029004801561041f5780601f106103fd57610100808354040283529182019161041f565b820191906000526020600020905b81548152906001019060200180831161040b575b50509150509081526020016040518091039020805490509050848111156104ff5780945060008381548110151561045257fe5b90600052602060002090600202016000018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104f75780601f106104cc576101008083540402835291602001916104f7565b820191906000526020600020905b8154815290600101906020018083116104da57829003601f168201915b505050505093505b50505b8080600101915050610206565b5060006010826040518082805190602001908083835b60208310151561054a5780518252602082019150602081019050602083039250610525565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020805490509050600260019054906101000a900460ff1660ff168110156106665760008190505b600260019054906101000a900460ff1660ff16811015610649576004816008811015156105cc57fe5b0154600f6000600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254019250508190555080806001019150506105a3565b5080600260016101000a81548160ff021916908360ff1602179055505b60008090505b600260019054906101000a900460ff1660ff168110156107c25760006010846040518082805190602001908083835b6020831015156106c0578051825260208201915060208101905060208303925061069b565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208281548110151561070057fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff166108fc600f60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460048560088110151561079657fe5b0154019081150290604051600060405180830381858888f193505050505050808060010191505061066c565b600260019054906101000a900460ff1660ff1690505b8181101561090a5760006010846040518082805190602001908083835b60208310151561081a57805182526020820191506020810190506020830392506107f5565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208281548110151561085a57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff166108fc600f60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549081150290604051600060405180830381858888f19350505050505080806001019150506107d8565b600e6000600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515610a3957600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600f6000600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549081150290604051600060405180830381858888f19350505050505b7f51e52a31977548a3c67fc31dcc97743791f4383a311fcc5470d88a7be8f06e1083856040518080602001838152602001828103825284818151815260200191508051906020019080838360005b83811015610aa2578082015181840152602081019050610a87565b50505050905090810190601f168015610acf5780820380516001836020036101000a031916815260200191505b50935050505060405180910390a1505050505b565b60643410151515610b5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f7061727469636970616e742066756e64206e6f7420656e6f756768000000000081525060200191505060405180910390fd5b6003544210151515610bd7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260098152602001807f707265706172696e67000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6001600260009054906101000a900460ff1660ff16141515610c61576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f6e6f742072756e6e696e6700000000000000000000000000000000000000000081525060200191505060405180910390fd5b600e60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615610cb857610f4d565b6000600f60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141580610d545750600d60009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15610d6757610d6233610f50565b610f4d565b600160008154809291906001019190505550600060408051908101604052808381526020013373ffffffffffffffffffffffffffffffffffffffff16815250908060018154018082558091505090600182039060005260206000209060020201600090919290919091506000820151816000019080519060200190610ded929190611066565b5060208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505034600f60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055507f13e1e749fcc13498ef55703026af9f512d4196e4a169c1dcfe303761370cef423382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610f11578082015181840152602081019050610ef6565b50505050905090810190601f168015610f3e5780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15b50565b600e60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615610fa757611063565b6001600e60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055507fb6d177e9c1090e879e384382951a234d03701128e360c4f2e0d568d51327b47881604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15b50565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106110a757805160ff19168380011785556110d5565b828001600101855582156110d5579182015b828111156110d45782518255916020019190600101906110b9565b5b5090506110e291906110e6565b5090565b61110891905b808211156111045760008160009055506001016110ec565b5090565b9056fea165627a7a7230582099def2d6bc1e1f8e66f336ade59c7d67939b29039b66b6a2d0b666eff4a799970029`

// DeployCalculationProc deploys a new Ethereum contract, binding an instance of CalculationProc to it.
func DeployCalculationProc(auth *bind.TransactOpts, backend bind.ContractBackend, id string, st *big.Int, fund *big.Int, participantCount uint8, distribute [8]*big.Int) (common.Address, *types.Transaction, *CalculationProc, error) {
	parsed, err := abi.JSON(strings.NewReader(CalculationProcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CalculationProcBin), backend, id, st, fund, participantCount, distribute)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CalculationProc{CalculationProcCaller: CalculationProcCaller{contract: contract}, CalculationProcTransactor: CalculationProcTransactor{contract: contract}, CalculationProcFilterer: CalculationProcFilterer{contract: contract}}, nil
}

// CalculationProc is an auto generated Go binding around an Ethereum contract.
type CalculationProc struct {
	CalculationProcCaller     // Read-only binding to the contract
	CalculationProcTransactor // Write-only binding to the contract
	CalculationProcFilterer   // Log filterer for contract events
}

// CalculationProcCaller is an auto generated read-only Go binding around an Ethereum contract.
type CalculationProcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalculationProcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CalculationProcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalculationProcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CalculationProcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CalculationProcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CalculationProcSession struct {
	Contract     *CalculationProc  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CalculationProcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CalculationProcCallerSession struct {
	Contract *CalculationProcCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CalculationProcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CalculationProcTransactorSession struct {
	Contract     *CalculationProcTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CalculationProcRaw is an auto generated low-level Go binding around an Ethereum contract.
type CalculationProcRaw struct {
	Contract *CalculationProc // Generic contract binding to access the raw methods on
}

// CalculationProcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CalculationProcCallerRaw struct {
	Contract *CalculationProcCaller // Generic read-only contract binding to access the raw methods on
}

// CalculationProcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CalculationProcTransactorRaw struct {
	Contract *CalculationProcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCalculationProc creates a new instance of CalculationProc, bound to a specific deployed contract.
func NewCalculationProc(address common.Address, backend bind.ContractBackend) (*CalculationProc, error) {
	contract, err := bindCalculationProc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CalculationProc{CalculationProcCaller: CalculationProcCaller{contract: contract}, CalculationProcTransactor: CalculationProcTransactor{contract: contract}, CalculationProcFilterer: CalculationProcFilterer{contract: contract}}, nil
}

// NewCalculationProcCaller creates a new read-only instance of CalculationProc, bound to a specific deployed contract.
func NewCalculationProcCaller(address common.Address, caller bind.ContractCaller) (*CalculationProcCaller, error) {
	contract, err := bindCalculationProc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CalculationProcCaller{contract: contract}, nil
}

// NewCalculationProcTransactor creates a new write-only instance of CalculationProc, bound to a specific deployed contract.
func NewCalculationProcTransactor(address common.Address, transactor bind.ContractTransactor) (*CalculationProcTransactor, error) {
	contract, err := bindCalculationProc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CalculationProcTransactor{contract: contract}, nil
}

// NewCalculationProcFilterer creates a new log filterer instance of CalculationProc, bound to a specific deployed contract.
func NewCalculationProcFilterer(address common.Address, filterer bind.ContractFilterer) (*CalculationProcFilterer, error) {
	contract, err := bindCalculationProc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CalculationProcFilterer{contract: contract}, nil
}

// bindCalculationProc binds a generic wrapper to an already deployed contract.
func bindCalculationProc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CalculationProcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CalculationProc *CalculationProcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CalculationProc.Contract.CalculationProcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CalculationProc *CalculationProcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CalculationProc.Contract.CalculationProcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CalculationProc *CalculationProcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CalculationProc.Contract.CalculationProcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CalculationProc *CalculationProcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CalculationProc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CalculationProc *CalculationProcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CalculationProc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CalculationProc *CalculationProcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CalculationProc.Contract.contract.Transact(opts, method, params...)
}

// Commit is a paid mutator transaction binding the contract method 0x9867db74.
//
// Solidity: function commit(string answerHash) returns()
func (_CalculationProc *CalculationProcTransactor) Commit(opts *bind.TransactOpts, answerHash string) (*types.Transaction, error) {
	return _CalculationProc.contract.Transact(opts, "commit", answerHash)
}

// Commit is a paid mutator transaction binding the contract method 0x9867db74.
//
// Solidity: function commit(string answerHash) returns()
func (_CalculationProc *CalculationProcSession) Commit(answerHash string) (*types.Transaction, error) {
	return _CalculationProc.Contract.Commit(&_CalculationProc.TransactOpts, answerHash)
}

// Commit is a paid mutator transaction binding the contract method 0x9867db74.
//
// Solidity: function commit(string answerHash) returns()
func (_CalculationProc *CalculationProcTransactorSession) Commit(answerHash string) (*types.Transaction, error) {
	return _CalculationProc.Contract.Commit(&_CalculationProc.TransactOpts, answerHash)
}

// Terminate is a paid mutator transaction binding the contract method 0x0c08bf88.
//
// Solidity: function terminate() returns()
func (_CalculationProc *CalculationProcTransactor) Terminate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CalculationProc.contract.Transact(opts, "terminate")
}

// Terminate is a paid mutator transaction binding the contract method 0x0c08bf88.
//
// Solidity: function terminate() returns()
func (_CalculationProc *CalculationProcSession) Terminate() (*types.Transaction, error) {
	return _CalculationProc.Contract.Terminate(&_CalculationProc.TransactOpts)
}

// Terminate is a paid mutator transaction binding the contract method 0x0c08bf88.
//
// Solidity: function terminate() returns()
func (_CalculationProc *CalculationProcTransactorSession) Terminate() (*types.Transaction, error) {
	return _CalculationProc.Contract.Terminate(&_CalculationProc.TransactOpts)
}

// CalculationProcCommittedIterator is returned from FilterCommitted and is used to iterate over the raw logs and unpacked data for Committed events raised by the CalculationProc contract.
type CalculationProcCommittedIterator struct {
	Event *CalculationProcCommitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CalculationProcCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CalculationProcCommitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CalculationProcCommitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CalculationProcCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CalculationProcCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CalculationProcCommitted represents a Committed event raised by the CalculationProc contract.
type CalculationProcCommitted struct {
	Participant common.Address
	AnsHash     string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCommitted is a free log retrieval operation binding the contract event 0x13e1e749fcc13498ef55703026af9f512d4196e4a169c1dcfe303761370cef42.
//
// Solidity: event committed(address participant, string ansHash)
func (_CalculationProc *CalculationProcFilterer) FilterCommitted(opts *bind.FilterOpts) (*CalculationProcCommittedIterator, error) {

	logs, sub, err := _CalculationProc.contract.FilterLogs(opts, "committed")
	if err != nil {
		return nil, err
	}
	return &CalculationProcCommittedIterator{contract: _CalculationProc.contract, event: "committed", logs: logs, sub: sub}, nil
}

// WatchCommitted is a free log subscription operation binding the contract event 0x13e1e749fcc13498ef55703026af9f512d4196e4a169c1dcfe303761370cef42.
//
// Solidity: event committed(address participant, string ansHash)
func (_CalculationProc *CalculationProcFilterer) WatchCommitted(opts *bind.WatchOpts, sink chan<- *CalculationProcCommitted) (event.Subscription, error) {

	logs, sub, err := _CalculationProc.contract.WatchLogs(opts, "committed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CalculationProcCommitted)
				if err := _CalculationProc.contract.UnpackLog(event, "committed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCommitted is a log parse operation binding the contract event 0x13e1e749fcc13498ef55703026af9f512d4196e4a169c1dcfe303761370cef42.
//
// Solidity: event committed(address participant, string ansHash)
func (_CalculationProc *CalculationProcFilterer) ParseCommitted(log types.Log) (*CalculationProcCommitted, error) {
	event := new(CalculationProcCommitted)
	if err := _CalculationProc.contract.UnpackLog(event, "committed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CalculationProcPunishedIterator is returned from FilterPunished and is used to iterate over the raw logs and unpacked data for Punished events raised by the CalculationProc contract.
type CalculationProcPunishedIterator struct {
	Event *CalculationProcPunished // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CalculationProcPunishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CalculationProcPunished)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CalculationProcPunished)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CalculationProcPunishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CalculationProcPunishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CalculationProcPunished represents a Punished event raised by the CalculationProc contract.
type CalculationProcPunished struct {
	Participant common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPunished is a free log retrieval operation binding the contract event 0xb6d177e9c1090e879e384382951a234d03701128e360c4f2e0d568d51327b478.
//
// Solidity: event punished(address participant)
func (_CalculationProc *CalculationProcFilterer) FilterPunished(opts *bind.FilterOpts) (*CalculationProcPunishedIterator, error) {

	logs, sub, err := _CalculationProc.contract.FilterLogs(opts, "punished")
	if err != nil {
		return nil, err
	}
	return &CalculationProcPunishedIterator{contract: _CalculationProc.contract, event: "punished", logs: logs, sub: sub}, nil
}

// WatchPunished is a free log subscription operation binding the contract event 0xb6d177e9c1090e879e384382951a234d03701128e360c4f2e0d568d51327b478.
//
// Solidity: event punished(address participant)
func (_CalculationProc *CalculationProcFilterer) WatchPunished(opts *bind.WatchOpts, sink chan<- *CalculationProcPunished) (event.Subscription, error) {

	logs, sub, err := _CalculationProc.contract.WatchLogs(opts, "punished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CalculationProcPunished)
				if err := _CalculationProc.contract.UnpackLog(event, "punished", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePunished is a log parse operation binding the contract event 0xb6d177e9c1090e879e384382951a234d03701128e360c4f2e0d568d51327b478.
//
// Solidity: event punished(address participant)
func (_CalculationProc *CalculationProcFilterer) ParsePunished(log types.Log) (*CalculationProcPunished, error) {
	event := new(CalculationProcPunished)
	if err := _CalculationProc.contract.UnpackLog(event, "punished", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CalculationProcTerminatedIterator is returned from FilterTerminated and is used to iterate over the raw logs and unpacked data for Terminated events raised by the CalculationProc contract.
type CalculationProcTerminatedIterator struct {
	Event *CalculationProcTerminated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CalculationProcTerminatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CalculationProcTerminated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CalculationProcTerminated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CalculationProcTerminatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CalculationProcTerminatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CalculationProcTerminated represents a Terminated event raised by the CalculationProc contract.
type CalculationProcTerminated struct {
	Ans string
	Cnt *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTerminated is a free log retrieval operation binding the contract event 0x51e52a31977548a3c67fc31dcc97743791f4383a311fcc5470d88a7be8f06e10.
//
// Solidity: event terminated(string ans, uint256 cnt)
func (_CalculationProc *CalculationProcFilterer) FilterTerminated(opts *bind.FilterOpts) (*CalculationProcTerminatedIterator, error) {

	logs, sub, err := _CalculationProc.contract.FilterLogs(opts, "terminated")
	if err != nil {
		return nil, err
	}
	return &CalculationProcTerminatedIterator{contract: _CalculationProc.contract, event: "terminated", logs: logs, sub: sub}, nil
}

// WatchTerminated is a free log subscription operation binding the contract event 0x51e52a31977548a3c67fc31dcc97743791f4383a311fcc5470d88a7be8f06e10.
//
// Solidity: event terminated(string ans, uint256 cnt)
func (_CalculationProc *CalculationProcFilterer) WatchTerminated(opts *bind.WatchOpts, sink chan<- *CalculationProcTerminated) (event.Subscription, error) {

	logs, sub, err := _CalculationProc.contract.WatchLogs(opts, "terminated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CalculationProcTerminated)
				if err := _CalculationProc.contract.UnpackLog(event, "terminated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTerminated is a log parse operation binding the contract event 0x51e52a31977548a3c67fc31dcc97743791f4383a311fcc5470d88a7be8f06e10.
//
// Solidity: event terminated(string ans, uint256 cnt)
func (_CalculationProc *CalculationProcFilterer) ParseTerminated(log types.Log) (*CalculationProcTerminated, error) {
	event := new(CalculationProcTerminated)
	if err := _CalculationProc.contract.UnpackLog(event, "terminated", log); err != nil {
		return nil, err
	}
	return event, nil
}
