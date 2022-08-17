package client

import (
	"github.com/beowulf-foundation/beowulf-go/transactions"
	"github.com/beowulf-foundation/beowulf-go/types"
)

//SetAsset returns data of type Asset
func SetAsset(amount float64, symbol string) *types.Asset {
	return &types.Asset{Amount: amount, Symbol: symbol}
}

//JSONTrxString generate Trx to String
func JSONTrxString(v *transactions.SignedTransaction) (string, error) {
	ans, err := types.JSONMarshal(v)
	if err != nil {
		return "", err
	}
	return string(ans), nil
}

//JSONOpString generate Operations to String
func JSONOpString(v []types.Operation) (string, error) {
	var tx types.Operations
	tx = append(tx, v...)
	ans, err := types.JSONMarshal(tx)
	if err != nil {
		return "", err
	}
	return string(ans), nil
}
