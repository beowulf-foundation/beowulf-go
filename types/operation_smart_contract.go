package types

import (
	"beowulf-go/encoding/transaction"
)

//SmartContractOperation represents transfer operation data.
type SmartContractOperation struct {
	RequiredOwners StringSlice `json:"required_owners"`
	Scid           string      `json:"scid"`
	ScOperation    string      `json:"sc_operation"`
	Fee            string      `json:"fee"`
}

//Type function that defines the type of operation SmartContractOperation.
func (op *SmartContractOperation) Type() OpType {
	return TypeSmartContract
}

//Data returns the operation data SmartContractOperation.
func (op *SmartContractOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type SmtCreateOperation to bytes.
func (op *SmartContractOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeSmartContract.Code()))
	//enc.Encode(op.RequiredOwners)
	// encode AccountAuths as map[string]uint16
	enc.EncodeUVarint(uint64(len(op.RequiredOwners)))
	for _, v := range op.RequiredOwners {
		enc.EncodeString(v)
	}
	enc.Encode(op.Scid)
	enc.Encode(op.ScOperation)
	enc.EncodeMoney(op.Fee)
	//enc.Encode(op.Extensions)
	//enc.EncodeUVarint(0)
	return enc.Err()
}
