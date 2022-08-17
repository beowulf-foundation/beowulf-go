package types

import (
	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
)

//SmtCreateOperation represents transfer operation data.
type SmtCreateOperation struct {
	ControlAccount string          `json:"control_account"`
	Symbol         *AssetSymbol    `json:"symbol"`
	Creator        string          `json:"creator"`
	SmtCreationFee string          `json:"smt_creation_fee"`
	Precision      uint8           `json:"precision"`
	Extensions     [][]interface{} `json:"extensions"`
	MaxSupply      uint64          `json:"max_supply"`
}

//Type function that defines the type of operation SmtCreateOperation.
func (op *SmtCreateOperation) Type() OpType {
	return TypeSmtCreate
}

//Data returns the operation data SmtCreateOperation.
func (op *SmtCreateOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type SmtCreateOperation to bytes.
func (op *SmtCreateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeSmtCreate.Code()))
	enc.Encode(op.ControlAccount)
	enc.Encode(op.Symbol)
	enc.Encode(op.Creator)
	enc.EncodeMoney(op.SmtCreationFee)
	enc.Encode(op.Precision)
	//enc.Encode(op.Extensions)
	enc.EncodeUVarint(0)
	enc.Encode(op.MaxSupply)
	return enc.Err()
}
