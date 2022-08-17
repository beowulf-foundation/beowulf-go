package types

import (
	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
)

//AccountSupernodeVoteOperation represents account_supernode_vote operation data.
type AccountSupernodeVoteOperation struct {
	Account   string `json:"account"`
	Supernode string `json:"supernode"`
	Approve   bool   `json:"approve"`
	Votes     int64  `json:"votes"`
	Fee       string `json:"fee"`
}

//Type function that defines the type of operation AccountSupernodeVoteOperation.
func (op *AccountSupernodeVoteOperation) Type() OpType {
	return TypeAccountSupernodeVote
}

//Data returns the operation data AccountSupernodeVoteOperation.
func (op *AccountSupernodeVoteOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type AccountSupernodeVoteOperation to bytes.
func (op *AccountSupernodeVoteOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeAccountSupernodeVote.Code()))
	enc.Encode(op.Account)
	enc.Encode(op.Supernode)
	enc.EncodeBool(op.Approve)
	enc.Encode(op.Votes)
	enc.EncodeMoney(op.Fee)
	return enc.Err()
}
