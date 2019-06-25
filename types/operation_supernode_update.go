package types

import (
	"beowulf-go/encoding/transaction"
)

//SupernodeUpdateOperation represents supernode_update operation data.
type SupernodeUpdateOperation struct {
	Owner           string           `json:"owner"`
	BlockSigningKey string           `json:"block_signing_key"`
	Fee             string           `json:"fee"`
}

//Type function that defines the type of operation SupernodeUpdateOperation.
func (op *SupernodeUpdateOperation) Type() OpType {
	return TypeSupernodeUpdate
}

//Data returns the operation data SupernodeUpdateOperation.
func (op *SupernodeUpdateOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type SupernodeUpdateOperation to bytes.
func (op *SupernodeUpdateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeSupernodeUpdate.Code()))
	enc.Encode(op.Owner)
	enc.EncodePubKey(op.BlockSigningKey)
	enc.EncodeMoney(op.Fee)
	return enc.Err()
}
