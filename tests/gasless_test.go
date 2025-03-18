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
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	gaslessImpl "github.com/kaiachain/kaia/kaiax/gasless/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

func TestGasless(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	config := params.MainnetChainConfig.Copy()
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Config = config

	fullNode, node, validator, chainId, workspace := newBlockchain(t, config, genesis)
	defer os.RemoveAll(workspace)
	defer func() {
		if err := fullNode.Stop(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}()

	// create account
	numAccounts := 1
	_, accounts, _ := createAccount(t, numAccounts, validator)

	// set unit price
	node.TxPool().SetGasPrice(common.Big1)

	var (
		optsOwner  = bind.NewKeyedTransactor(accounts[0].Keys[0])
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
	)

	// send approveTx
	approveTx := gaslessImpl.MakeApproveTx(t, accounts[0].Keys[0], accounts[0].Nonce, gaslessImpl.ApproveArgs{Spender: common.HexToAddress("0x1234"), Amount: big.NewInt(1000000)}, chainId)
	err := transactor.SendTransaction(optsOwner.Context, approveTx)
	if err != nil {
		t.Fatal(err)
	}
	accounts[0].Nonce += 1

	// send swapTx
	swapTx := gaslessImpl.MakeSwapTx(t, accounts[0].Keys[0], accounts[0].Nonce, gaslessImpl.SwapArgs{Token: common.HexToAddress("0xabcd"), AmountIn: big.NewInt(10), MinAmountOut: big.NewInt(100), AmountRepay: big.NewInt(2021000)}, chainId)
	err = transactor.SendTransaction(optsOwner.Context, swapTx)
	if err != nil {
		t.Fatal(err)
	}
	accounts[0].Nonce += 1

	chain := node.BlockChain().(*blockchain.BlockChain)

	approveTxReceipt := waitReceipt(chain, approveTx.Hash())
	require.NotNil(t, approveTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, approveTxReceipt.Status, "approveTx failed")

	swapTxReceipt := waitReceipt(chain, swapTx.Hash())
	require.NotNil(t, swapTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, swapTxReceipt.Status, "swapTx failed")
}
