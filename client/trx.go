package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"beowulf-go/api"
	"beowulf-go/transactions"
	"beowulf-go/types"
)

//BResp of response when sending a transaction.
type BResp struct {
	ID       string
	//BlockNum int32
	//TrxNum   int32
	//Expired  bool
	//CreatedTime int64
	JSONTrx  string
}

//OperResp type is returned when the operation is performed.
type OperResp struct {
	NameOper string
	Bresp    *BResp
}

//SendTrx generates and sends an array of transactions to BEOWULF.
func (client *Client) SendTrx(strx []types.Operation) (*BResp, error) {
	var bresp BResp

	// Getting the necessary parameters
	props, err := client.API.GetDynamicGlobalProperties()
	if err != nil {
		return nil, err
	}

	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return nil, err
	}
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
		Extensions: [][]interface{}{},
	})

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(time.Duration(1) * time.Hour).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	// Obtain the key required for signing
	privKeys, err := client.SigningKeys(strx[0])
	if err != nil {
		return nil, err
	}

	// Sign the transaction
	txId, err := tx.Sign(privKeys, client.chainID)
	if err != nil {
		return nil, err
	}

	// Sending a transaction
	var resp *api.AsyncBroadcastResponse
	resp, err = client.API.BroadcastTransaction(tx.Transaction)
	//if client.AsyncProtocol {
	//	resp, err = client.API.BroadcastTransaction(tx.Transaction)
	//} else {
	//	resp, err = client.API.BroadcastTransactionSynchronous(tx.Transaction)
	//}

	bresp.JSONTrx, _ = JSONTrxString(tx)

	if err != nil {
		return &bresp, err
	}
	//if resp != nil && !client.AsyncProtocol {
	if resp != nil{
		txIdRes, _ := hex.DecodeString(resp.ID)
		fmt.Println(txIdRes)
		if(txId != resp.ID){
			return nil, errors.New("TransactionID is not mapped")
		}
		bresp.ID = resp.ID
		//bresp.BlockNum = resp.BlockNum
		//bresp.TrxNum = resp.TrxNum
		//bresp.Expired = resp.Expired
		//bresp.CreatedTime = resp.CreatedTime

		return &bresp, nil
	}

	return &bresp, nil
}

func (client *Client) GetTrx(strx []types.Operation) (*types.Transaction, error) {
	// Getting the necessary parameters
	props, err := client.API.GetDynamicGlobalProperties()
	if err != nil {
		return nil, err
	}

	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
		Extensions: [][]interface{}{},
	}

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(1 * time.Second).UTC()
	//expTime := time.Now().Add(59 * time.Minute).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	return tx, nil
}
