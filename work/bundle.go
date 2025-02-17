package work

import (
	"github.com/kaiachain/kaia/blockchain/types"
)

type bundle struct {
	BundleTxs types.Transactions // each element can be either *types.Transaction, or TxGenerator
	// TargetTxHash common.Hash   // BundleTxs is placed AFTER the target tx. If empty hash, it is placed at the very front.
}

func arrayfy(txs types.TransactionsByPriceAndNonce) types.Transactions {
	arrayTxs := types.Transactions{}
	for {
		tx := txs.Peek()
		if tx == nil {
			break
		}
		arrayTxs = append(arrayTxs, tx)

		txs.Shift()
	}
	return arrayTxs
}
