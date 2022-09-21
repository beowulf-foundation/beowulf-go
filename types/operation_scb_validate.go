package types

import (
	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
)

//SmartContractOperation represents transfer operation data.
type ScbValidateOperation struct {
	Committer   string `json:"committer"`
	Supernode   string `json:"supernode"`
	Scid        string `json:"scid"`
	ScOperation string `json:"sc_operation"`
	TimeCommit  UInt32 `json:"time_commit"`
	Fee         string `json:"fee"`
}

//Type function that defines the type of operation SmartContractOperation.
func (op *ScbValidateOperation) Type() OpType {
	return TypeScbValidate
}

//Data returns the operation data SmartContractOperation.
func (op *ScbValidateOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type SmtCreateOperation to bytes.
func (op *ScbValidateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeScbValidate.Code()))
	//enc.Encode(op.RequiredOwners)
	// encode AccountAuths as map[string]uint16
	enc.EncodeString(op.Committer)
	enc.Encode(op.Supernode)
	enc.Encode(op.Scid)
	enc.Encode(op.ScOperation)
	enc.Encode(op.TimeCommit)
	enc.EncodeMoney(op.Fee)
	//enc.Encode(op.Extensions)
	//enc.EncodeUVarint(0)
	return enc.Err()
}
