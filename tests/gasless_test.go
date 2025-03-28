// Modifications Copyright 2025 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package tests

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	uniswapFactoryContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/factory"
	uniswapRouterContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/router"
	gaslessContract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip247"
	testingContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	testingGaslessContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts/gasless"

	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/node/cn"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

var (
	bigKaia = big.NewInt(params.KAIA)
	bigGkei = big.NewInt(params.Gkei)
)

func TestGasless(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlInfo)

	// prepare chain configuration
	config := params.MainnetChainConfig.Copy()
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = nil
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.PragueCompatibleBlock = big.NewInt(0)

	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)

	fullNode, node, validator, _, workspace := newBlockchain(t, config, nil)
	defer func() {
		os.RemoveAll(workspace)
		if err := fullNode.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	numAccounts := 1
	_, accounts, _ := createAccount(t, numAccounts, validator)

	var (
		owner = validator
	)

	// deploy contracts
	_, testTokenAddr, _, testTokenContract := deployTestToken(t, node, owner, owner.Addr)
	_, wkaiaAddr, _, wkaiaContract := deployWKAIA(t, node, owner)
	_, factoryAddr, _, factoryContract := deployUniswapV2Factory(t, node, owner, owner.Addr)
	_, routerAddr, _, routerContract := deployUniswapV2Router02(t, node, owner, factoryAddr, wkaiaAddr)
	_, gsrAddr, _, gsrContract := deployGaslessSwapRouter(t, node, owner, wkaiaAddr)

	// Register GaslessSwapRouter address in Registry
	/* ------------- process ------------- */

	// Set up initial liquidity
	contracts := contractsForGasless{
		testTokenAddr:     testTokenAddr,
		testTokenContract: testTokenContract,
		wkaiaAddr:         wkaiaAddr,
		wkaiaContract:     wkaiaContract,
		factoryAddr:       factoryAddr,
		factoryContract:   factoryContract,
		routerAddr:        routerAddr,
		routerContract:    routerContract,
		gsrAddr:           gsrAddr,
		gsrContract:       gsrContract,
	}
	setupLiquidity(t, owner, contracts, node)

	// update gasless module
	node.GetGaslessModule().UpdateRouter(gsrAddr)
	node.GetGaslessModule().UpdateAllowedToken(testTokenAddr)

	/* ------------------------------------ main test ------------------------------------- */
	swapAmmount := new(big.Int).Mul(big.NewInt(1), bigKaia)
	amountsOut, err := routerContract.GetAmountsOut(&bind.CallOpts{}, swapAmmount, []common.Address{testTokenAddr, wkaiaAddr})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(amountsOut)
	swapExpectedOutput := amountsOut[1]
	// swapExpectedOutput := big.NewInt(100)

	var (
		gasPriceBN    = new(big.Int).Mul(big.NewInt(50), bigGkei)
		R1            = new(big.Int).Mul(big.NewInt(21000), gasPriceBN)
		R2            = new(big.Int).Mul(big.NewInt(100000), gasPriceBN)
		R3            = new(big.Int).Mul(big.NewInt(500000), gasPriceBN)
		ammontRepay   = new(big.Int).Add(R1, new(big.Int).Add(R2, R3))
		transferToken = new(big.Int).Mul(big.NewInt(100), bigKaia)
		margin        = new(big.Int).Div(swapExpectedOutput, big.NewInt(100))
		minAmountOut  = new(big.Int).Add(ammontRepay, margin)
	)

	// transfer test token
	testTokenTransferTx, err := testTokenContract.Transfer(bind.NewKeyedTransactor(owner.Keys[0]), accounts[0].Addr, transferToken)
	if err != nil {
		t.Fatal(err)
	}
	testTokenTransferReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), testTokenTransferTx.Hash())
	if testTokenTransferReceipt == nil {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	balanceOfTestAcc, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, accounts[0].Addr)
	balanceOfOwner, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	fmt.Println("test acc balance is: ", balanceOfTestAcc)
	fmt.Println("owner balance is: ", balanceOfOwner)

	// send approveTx
	optsForApprove := bind.NewKeyedTransactor(accounts[0].Keys[0])
	optsForApprove.GasLimit = 300000
	optsForApprove.Nonce = big.NewInt(int64(accounts[0].Nonce))
	approveTx, err := testTokenContract.Approve(optsForApprove, gsrAddr, swapAmmount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("approveTxHash", approveTx.Hash().Hex())
	accounts[0].Nonce += 1

	// send swapTx
	optsForSwap := bind.NewKeyedTransactor(accounts[0].Keys[0])
	optsForSwap.GasLimit = 300000
	optsForSwap.Nonce = big.NewInt(int64(accounts[0].Nonce))
	swapTx, err := gsrContract.SwapForGas(optsForSwap, testTokenAddr, swapAmmount, minAmountOut, ammontRepay)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("swapTxHash", swapTx.Hash().Hex(), swapTx)
	accounts[0].Nonce += 1

	pending, queue := node.TxPool().Stats()
	t.Log("pending", pending, "queue", queue)

	chain := node.BlockChain().(*blockchain.BlockChain)

	approveTxReceipt := waitReceipt(chain, approveTx.Hash())
	require.NotNil(t, approveTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, approveTxReceipt.Status, "approveTx failed")

	swapTxReceipt := waitReceipt(chain, swapTx.Hash())
	require.NotNil(t, swapTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, swapTxReceipt.Status, "swapTx failed")
	/* --------------------------------------------------------------------------- */
}

func deployTestToken(t *testing.T, node *cn.CN, owner *TestAccountType, initialHolder common.Address,
) (uint64, common.Address, *types.Transaction, *testingGaslessContracts.TestToken) {
	var (
		optsOwner  = bind.NewKeyedTransactor(owner.Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// Deploy contract
	addr, tx, contract, err := testingGaslessContracts.DeployTestToken(optsOwner, transactor, initialHolder)
	if err != nil {
		t.Fatal(err)
	}

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("TestToken deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return num, addr, tx, contract
}

func deployWKAIA(t *testing.T, node *cn.CN, owner *TestAccountType,
) (uint64, common.Address, *types.Transaction, *testingContracts.WKAIA) {
	var (
		optsOwner  = bind.NewKeyedTransactor(owner.Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// Deploy contract
	addr, tx, contract, err := testingContracts.DeployWKAIA(optsOwner, transactor)
	if err != nil {
		t.Fatal(err)
	}

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("WKAIA deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return num, addr, tx, contract
}

func deployUniswapV2Factory(t *testing.T, node *cn.CN, owner *TestAccountType, feeToSetter common.Address,
) (uint64, common.Address, *types.Transaction, *uniswapFactoryContracts.UniswapV2Factory) {
	var (
		optsOwner  = bind.NewKeyedTransactor(owner.Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// Deploy contract
	addr, tx, contract, err := uniswapFactoryContracts.DeployUniswapV2Factory(optsOwner, transactor, feeToSetter)
	if err != nil {
		t.Fatal(err)
	}

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("UniswapV2Factory deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return num, addr, tx, contract
}

func deployUniswapV2Router02(t *testing.T, node *cn.CN, owner *TestAccountType, factory, wkaia common.Address,
) (uint64, common.Address, *types.Transaction, *uniswapRouterContracts.UniswapV2Router02) {
	var (
		optsOwner  = bind.NewKeyedTransactor(owner.Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// Deploy contract
	addr, tx, contract, err := uniswapRouterContracts.DeployUniswapV2Router02(optsOwner, transactor, factory, wkaia)
	if err != nil {
		t.Fatal(err)
	}

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("UniswapV2Router02 deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return num, addr, tx, contract
}

func deployGaslessSwapRouter(t *testing.T, node *cn.CN, owner *TestAccountType, wkaiaAddr common.Address,
) (uint64, common.Address, *types.Transaction, *gaslessContract.GaslessSwapRouter) {
	var (
		optsOwner  = bind.NewKeyedTransactor(owner.Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// Deploy contract
	addr, tx, contract, err := gaslessContract.DeployGaslessSwapRouter(optsOwner, transactor, wkaiaAddr)
	if err != nil {
		t.Fatal(err)
	}

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GaslessSwapRouter deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return num, addr, tx, contract
}

type contractsForGasless struct {
	testTokenAddr     common.Address
	testTokenContract *testingGaslessContracts.TestToken
	wkaiaAddr         common.Address
	wkaiaContract     *testingContracts.WKAIA
	factoryAddr       common.Address
	factoryContract   *uniswapFactoryContracts.UniswapV2Factory
	routerAddr        common.Address
	routerContract    *uniswapRouterContracts.UniswapV2Router02
	gsrAddr           common.Address
	gsrContract       *gaslessContract.GaslessSwapRouter
}

func setupLiquidity(t *testing.T, owner *TestAccountType, contracts contractsForGasless, node *cn.CN) {
	var (
		testTokenAddr     = contracts.testTokenAddr
		testTokenContract = contracts.testTokenContract
		wkaiaAddr         = contracts.wkaiaAddr
		wkaiaContract     = contracts.wkaiaContract
		factoryAddr       = contracts.factoryAddr
		factoryContract   = contracts.factoryContract
		routerAddr        = contracts.routerAddr
		routerContract    = contracts.routerContract
		gsrContract       = contracts.gsrContract
		initialLiquidity  = new(big.Int).Mul(big.NewInt(1000), bigKaia)
	)

	/* ------------- create pair ------------- */
	createPairTx, err := factoryContract.CreatePair(bind.NewKeyedTransactor(owner.Keys[0]), testTokenAddr, wkaiaAddr)
	if err != nil {
		t.Fatal(err)
	}
	createPairReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), createPairTx.Hash())
	if createPairReceipt == nil || createPairReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	pairAddr, _ := factoryContract.GetPair(&bind.CallOpts{}, testTokenAddr, wkaiaAddr)

	fmt.Println("---------------addresses----------------")
	fmt.Println("testTokenAddr", testTokenAddr.Hex())
	fmt.Println("wkaiaAddr", wkaiaAddr.Hex())
	fmt.Println("factoryAddr", factoryAddr.Hex())
	fmt.Println("routerAddr", routerAddr.Hex())
	fmt.Println("pairAddr", pairAddr.Hex())
	fmt.Println("-------------------------------")

	/* ------------- deposit ------------- */
	optsForDeposit := bind.NewKeyedTransactor(owner.Keys[0])
	optsForDeposit.Value = initialLiquidity
	optsForDeposit.GasLimit = 300000
	depositTx, err := wkaiaContract.Deposit(optsForDeposit)
	if err != nil {
		t.Fatal(err)
	}
	depositReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), depositTx.Hash())
	if depositReceipt == nil || depositReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	/* ------------- approve(test token) ------------- */
	testTokenApproveTx, err := testTokenContract.Approve(bind.NewKeyedTransactor(owner.Keys[0]), routerAddr, initialLiquidity)
	if err != nil {
		t.Fatal(err)
	}
	testTokenApproveReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), testTokenApproveTx.Hash())
	if testTokenApproveReceipt == nil || testTokenApproveReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	/* ------------- approve(wkaia) ------------- */
	wkaiaApproveTx, err := wkaiaContract.Approve(bind.NewKeyedTransactor(owner.Keys[0]), routerAddr, initialLiquidity)
	if err != nil {
		t.Fatal(err)
	}
	wkaiaApproveReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), wkaiaApproveTx.Hash())
	if wkaiaApproveReceipt == nil || wkaiaApproveReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	balanceOfWKAIA, _ := wkaiaContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	balanceOfTestToken, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	wallowance, _ := wkaiaContract.Allowance(&bind.CallOpts{}, owner.Addr, routerAddr)
	tallowance, _ := testTokenContract.Allowance(&bind.CallOpts{}, owner.Addr, routerAddr)
	t.Log(balanceOfWKAIA, balanceOfTestToken, "balance of tokens")
	t.Log(wallowance, tallowance, "allowances of tokens")

	/* ------------- add liquidity ------------- */
	optsForAddLiquidity := bind.NewKeyedTransactor(owner.Keys[0])
	optsForAddLiquidity.GasLimit = 3000000
	deadline := time.Now().Unix() + 60*20
	addLiquidityTx, err := routerContract.AddLiquidity(optsForAddLiquidity, testTokenAddr, wkaiaAddr,
		initialLiquidity, initialLiquidity, common.Big0, common.Big0, owner.Addr, big.NewInt(deadline))
	if err != nil {
		t.Fatal(err)
	}
	addLiquidityReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), addLiquidityTx.Hash())
	if addLiquidityReceipt == nil || addLiquidityReceipt.Status != types.ReceiptStatusSuccessful {
		t.Log(addLiquidityReceipt)
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	// Add token to gsr
	optsForAddToken := bind.NewKeyedTransactor(owner.Keys[0])
	optsForAddToken.GasLimit = 300000
	addTokenTx, err := gsrContract.AddToken(optsForAddToken, testTokenAddr, factoryAddr, routerAddr)
	if err != nil {
		t.Fatal(err)
	}
	addTokenReceipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), addTokenTx.Hash())
	if addTokenReceipt == nil || addTokenReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("timeout")
	}
	owner.Nonce += 1

	var (
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)
	pairContract, err := uniswapRouterContracts.NewIUniswapV2Pair(pairAddr, transactor)
	if err != nil {
		t.Fatal(err)
	}
	reserves, err := pairContract.GetReserves(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(reserves, "----------------------------reserves---------------")
}
