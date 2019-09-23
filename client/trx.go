package client

import (
	"beowulf-go/api"
	"beowulf-go/config"
	"beowulf-go/transactions"
	"beowulf-go/types"
	"time"
)

var RefBlockMap = make(map[time.Time]uint32)

//BResp of response when sending a transaction.
type BResp struct {
	ID string
	JSONTrx string
}

//OperResp type is returned when the operation is performed.
type OperResp struct {
	NameOper string
	Bresp    *BResp
}

//Get HeadBlockNumber from mem before getting from Blockchain
func (client *Client) GetHeadBlockNum() (uint32, error) {
	if len(RefBlockMap) > 0 {
		for k := range RefBlockMap {
			old := k.Add(config.GET_HEAD_BLOCK_NUM_TIMEOUT_IN_MIN * time.Minute)
			now := time.Now().UTC()
			if old.Before(now) {
				delete(RefBlockMap, k)
				props, err := client.API.GetDynamicGlobalProperties()
				if err != nil {
					return 0, err
				}
				refBlockNum := props.HeadBlockNumber
				if refBlockNum > config.HEAD_BLOCK_NUM_SPAN {
					refBlockNum -= config.HEAD_BLOCK_NUM_SPAN
				}
				RefBlockMap[now] = refBlockNum
				return refBlockNum, nil
			}
			return RefBlockMap[k], nil
		}
	}
	props, err := client.API.GetDynamicGlobalProperties()
	if err != nil {
		return 0, err
	}
	refBlockNum := props.HeadBlockNumber
	if refBlockNum > config.HEAD_BLOCK_NUM_SPAN {
		refBlockNum -= config.HEAD_BLOCK_NUM_SPAN
	}
	now := time.Now().UTC()
	RefBlockMap[now] = refBlockNum
	return refBlockNum, nil
}

//SendTrx generates and sends an array of transactions to BEOWULF.
func (client *Client) SendTrx(strx []types.Operation) (*BResp, error) {
	var bresp BResp

	// Getting the necessary parameters
	refBlockNum, err := client.GetHeadBlockNum()
	if err != nil {
		return nil, err
	}
	block, err := client.API.GetBlock(refBlockNum)
	if err != nil {
		return nil, err
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return nil, err
	}
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(refBlockNum),
		RefBlockPrefix: refBlockPrefix,
		Extensions:     [][]interface{}{},
	})

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
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
	tx.Transaction.Signatures = []string{}
	txId, err := tx.Sign(privKeys, client.chainID)
	if err != nil || txId == "" {
		return nil, err
	}

	// Sending a transaction
	if client.AsyncProtocol {
		var resp *api.AsyncBroadcastResponse
		resp, err = client.API.BroadcastTransaction(tx.Transaction)
		if resp != nil {
			bresp.ID = resp.ID
		}
	} else {
		var resp *api.BroadcastResponse
		resp, err = client.API.BroadcastTransactionSynchronous(tx.Transaction)
		if resp != nil {
			bresp.ID = resp.ID
		}
	}

	bresp.JSONTrx, _ = JSONTrxString(tx)

	if err != nil {
		return &bresp, err
	}

	return &bresp, nil
}

func (client *Client) GetTrx(strx []types.Operation) (*types.Transaction, error) {
	// Getting the necessary parameters
	refBlockNum, err := client.GetHeadBlockNum()
	if err != nil {
		return nil, err
	}
	block, err := client.API.GetBlock(refBlockNum)
	if err != nil {
		return nil, err
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(refBlockNum),
		RefBlockPrefix: refBlockPrefix,
		Extensions:     [][]interface{}{},
	}

	// Adding Operations to a Transaction
	for _, val := range strx {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
	//expTime := time.Now().Add(59 * time.Minute).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	return tx, nil
}
