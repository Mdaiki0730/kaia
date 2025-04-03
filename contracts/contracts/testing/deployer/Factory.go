// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deployer

import (
	"errors"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = kaia.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FactoryMetaData contains all meta data concerning the Factory contract.
var FactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"name\":\"Deployed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"code\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"name\":\"deploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"9c4ae2d0": "deploy(bytes,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b506101a6806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80639c4ae2d014610030575b600080fd5b61004361003e3660046100bb565b610045565b005b6000818351602085016000f59050803b61005e57600080fd5b604080516001600160a01b0383168152602081018490527fb03c53b28e78a88e31607a27e1fa48234dce28d5d9d9ec7b295aeb02e674a1e1910160405180910390a1505050565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156100ce57600080fd5b823567ffffffffffffffff808211156100e657600080fd5b818501915085601f8301126100fa57600080fd5b81358181111561010c5761010c6100a5565b604051601f8201601f19908116603f01168101908382118183101715610134576101346100a5565b8160405282815288602084870101111561014d57600080fd5b82602086016020830137600060209382018401529896909101359650505050505056fea264697066735822122026aafebf064e9d81c9da61f3fe7ecbf59287d5a9d71c729e4ac88de771584adc64736f6c634300080d0033",
}

// FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use FactoryMetaData.ABI instead.
var FactoryABI = FactoryMetaData.ABI

// FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const FactoryBinRuntime = `608060405234801561001057600080fd5b506004361061002b5760003560e01c80639c4ae2d014610030575b600080fd5b61004361003e3660046100bb565b610045565b005b6000818351602085016000f59050803b61005e57600080fd5b604080516001600160a01b0383168152602081018490527fb03c53b28e78a88e31607a27e1fa48234dce28d5d9d9ec7b295aeb02e674a1e1910160405180910390a1505050565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156100ce57600080fd5b823567ffffffffffffffff808211156100e657600080fd5b818501915085601f8301126100fa57600080fd5b81358181111561010c5761010c6100a5565b604051601f8201601f19908116603f01168101908382118183101715610134576101346100a5565b8160405282815288602084870101111561014d57600080fd5b82602086016020830137600060209382018401529896909101359650505050505056fea264697066735822122026aafebf064e9d81c9da61f3fe7ecbf59287d5a9d71c729e4ac88de771584adc64736f6c634300080d0033`

// Deprecated: Use FactoryMetaData.Sigs instead.
// FactoryFuncSigs maps the 4-byte function signature to its string representation.
var FactoryFuncSigs = FactoryMetaData.Sigs

// FactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FactoryMetaData.Bin instead.
var FactoryBin = FactoryMetaData.Bin

// DeployFactory deploys a new Kaia contract, binding an instance of Factory to it.
func DeployFactory(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Factory, error) {
	parsed, err := FactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FactoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Factory{FactoryCaller: FactoryCaller{contract: contract}, FactoryTransactor: FactoryTransactor{contract: contract}, FactoryFilterer: FactoryFilterer{contract: contract}}, nil
}

// Factory is an auto generated Go binding around a Kaia contract.
type Factory struct {
	FactoryCaller     // Read-only binding to the contract
	FactoryTransactor // Write-only binding to the contract
	FactoryFilterer   // Log filterer for contract events
}

// FactoryCaller is an auto generated read-only Go binding around a Kaia contract.
type FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FactoryTransactor is an auto generated write-only Go binding around a Kaia contract.
type FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FactoryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FactorySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type FactorySession struct {
	Contract     *Factory          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FactoryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type FactoryCallerSession struct {
	Contract *FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// FactoryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type FactoryTransactorSession struct {
	Contract     *FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// FactoryRaw is an auto generated low-level Go binding around a Kaia contract.
type FactoryRaw struct {
	Contract *Factory // Generic contract binding to access the raw methods on
}

// FactoryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type FactoryCallerRaw struct {
	Contract *FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// FactoryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type FactoryTransactorRaw struct {
	Contract *FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFactory creates a new instance of Factory, bound to a specific deployed contract.
func NewFactory(address common.Address, backend bind.ContractBackend) (*Factory, error) {
	contract, err := bindFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Factory{FactoryCaller: FactoryCaller{contract: contract}, FactoryTransactor: FactoryTransactor{contract: contract}, FactoryFilterer: FactoryFilterer{contract: contract}}, nil
}

// NewFactoryCaller creates a new read-only instance of Factory, bound to a specific deployed contract.
func NewFactoryCaller(address common.Address, caller bind.ContractCaller) (*FactoryCaller, error) {
	contract, err := bindFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FactoryCaller{contract: contract}, nil
}

// NewFactoryTransactor creates a new write-only instance of Factory, bound to a specific deployed contract.
func NewFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*FactoryTransactor, error) {
	contract, err := bindFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FactoryTransactor{contract: contract}, nil
}

// NewFactoryFilterer creates a new log filterer instance of Factory, bound to a specific deployed contract.
func NewFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*FactoryFilterer, error) {
	contract, err := bindFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FactoryFilterer{contract: contract}, nil
}

// bindFactory binds a generic wrapper to an already deployed contract.
func bindFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Factory *FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Factory.Contract.FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Factory *FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Factory.Contract.FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Factory *FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Factory.Contract.FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Factory *FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Factory *FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Factory *FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Factory.Contract.contract.Transact(opts, method, params...)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes code, uint256 salt) returns()
func (_Factory *FactoryTransactor) Deploy(opts *bind.TransactOpts, code []byte, salt *big.Int) (*types.Transaction, error) {
	return _Factory.contract.Transact(opts, "deploy", code, salt)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes code, uint256 salt) returns()
func (_Factory *FactorySession) Deploy(code []byte, salt *big.Int) (*types.Transaction, error) {
	return _Factory.Contract.Deploy(&_Factory.TransactOpts, code, salt)
}

// Deploy is a paid mutator transaction binding the contract method 0x9c4ae2d0.
//
// Solidity: function deploy(bytes code, uint256 salt) returns()
func (_Factory *FactoryTransactorSession) Deploy(code []byte, salt *big.Int) (*types.Transaction, error) {
	return _Factory.Contract.Deploy(&_Factory.TransactOpts, code, salt)
}

// FactoryDeployedIterator is returned from FilterDeployed and is used to iterate over the raw logs and unpacked data for Deployed events raised by the Factory contract.
type FactoryDeployedIterator struct {
	Event *FactoryDeployed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FactoryDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FactoryDeployed)
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
		it.Event = new(FactoryDeployed)
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
func (it *FactoryDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FactoryDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FactoryDeployed represents a Deployed event raised by the Factory contract.
type FactoryDeployed struct {
	Addr common.Address
	Salt *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDeployed is a free log retrieval operation binding the contract event 0xb03c53b28e78a88e31607a27e1fa48234dce28d5d9d9ec7b295aeb02e674a1e1.
//
// Solidity: event Deployed(address addr, uint256 salt)
func (_Factory *FactoryFilterer) FilterDeployed(opts *bind.FilterOpts) (*FactoryDeployedIterator, error) {

	logs, sub, err := _Factory.contract.FilterLogs(opts, "Deployed")
	if err != nil {
		return nil, err
	}
	return &FactoryDeployedIterator{contract: _Factory.contract, event: "Deployed", logs: logs, sub: sub}, nil
}

// WatchDeployed is a free log subscription operation binding the contract event 0xb03c53b28e78a88e31607a27e1fa48234dce28d5d9d9ec7b295aeb02e674a1e1.
//
// Solidity: event Deployed(address addr, uint256 salt)
func (_Factory *FactoryFilterer) WatchDeployed(opts *bind.WatchOpts, sink chan<- *FactoryDeployed) (event.Subscription, error) {

	logs, sub, err := _Factory.contract.WatchLogs(opts, "Deployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FactoryDeployed)
				if err := _Factory.contract.UnpackLog(event, "Deployed", log); err != nil {
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

// ParseDeployed is a log parse operation binding the contract event 0xb03c53b28e78a88e31607a27e1fa48234dce28d5d9d9ec7b295aeb02e674a1e1.
//
// Solidity: event Deployed(address addr, uint256 salt)
func (_Factory *FactoryFilterer) ParseDeployed(log types.Log) (*FactoryDeployed, error) {
	event := new(FactoryDeployed)
	if err := _Factory.contract.UnpackLog(event, "Deployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
