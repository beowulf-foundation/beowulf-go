package types

import (
	"encoding/json"
)

//OperationObject type from parameter JSON
type OperationObject struct {
	TransactionID          string    `json:"trx_id"`
	BlockNumber            uint32    `json:"block"`
	TransactionInBlock     uint32    `json:"trx_in_block"`
	OperationInTransaction uint32    `json:"op_in_trx"`
	VirtualOperation       uint32    `json:"virtual_op"`
	Timestamp              *Time     `json:"timestamp"`
	Operation              Operation `json:"op"`
	OperationType          OpType    `json:"-"`
}

type rawOperationObject struct {
	TransactionID          string          `json:"trx_id"`
	BlockNumber            uint32          `json:"block"`
	TransactionInBlock     uint32          `json:"trx_in_block"`
	OperationInTransaction uint32          `json:"op_in_trx"`
	VirtualOperation       uint32          `json:"virtual_op"`
	Timestamp              *Time           `json:"timestamp"`
	Operation              *operationTuple `json:"op"`
}

//UnmarshalJSON unpacking the JSON parameter in the OperationObject type.
func (op *OperationObject) UnmarshalJSON(p []byte) error {
	var raw rawOperationObject
	if err := json.Unmarshal(p, &raw); err != nil {
		return err
	}

	op.TransactionID = raw.TransactionID
	op.BlockNumber = raw.BlockNumber
	op.TransactionInBlock = raw.TransactionInBlock
	op.OperationInTransaction = raw.OperationInTransaction
	op.VirtualOperation = raw.VirtualOperation
	op.Timestamp = raw.Timestamp
	op.Operation = raw.Operation.Data
	op.OperationType = raw.Operation.Type
	return nil
}

//MarshalJSON function for packing the OperationObject type in JSON.
func (op *OperationObject) MarshalJSON() ([]byte, error) {
	return JSONMarshal(&rawOperationObject{
		TransactionID:          op.TransactionID,
		BlockNumber:            op.BlockNumber,
		TransactionInBlock:     op.TransactionInBlock,
		OperationInTransaction: op.OperationInTransaction,
		VirtualOperation:       op.VirtualOperation,
		Timestamp:              op.Timestamp,
		Operation:              &operationTuple{op.Operation.Type(), op.Operation},
	})
}
