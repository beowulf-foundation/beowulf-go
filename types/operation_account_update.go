package types

import (
	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
)

//AccountUpdateOperation represents account_update operation data.
type AccountUpdateOperation struct {
	Account      string           `json:"account"`
	Owner        *Authority       `json:"owner,omitempty"`
	JSONMetadata *AccountMetadata `json:"json_metadata"`
	Fee          string           `json:"fee"`
}

//Type function that defines the type of operation AccountUpdateOperation.
func (op *AccountUpdateOperation) Type() OpType {
	return TypeAccountUpdate
}

//Data returns the operation data AccountUpdateOperation.
func (op *AccountUpdateOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type AccountUpdateOperation to bytes.
func (op *AccountUpdateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeAccountUpdate.Code()))
	enc.EncodeString(op.Account)
	if op.Owner != nil {
		enc.Encode(byte(1))
		enc.Encode(op.Owner)
	} else {
		enc.Encode(byte(0))
	}
	enc.Encode(op.JSONMetadata)
	enc.EncodeMoney(op.Fee)
	return enc.Err()
}
