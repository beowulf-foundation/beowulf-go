package types

import (
	"beowulf-go/encoding/transaction"
)

//TransferToVestingOperation represents transfer_to_vesting operation data.
type TransferToVestingOperation struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Fee    string `json:"fee"`
}

//Type function that defines the type of operation TransferToVestingOperation.
func (op *TransferToVestingOperation) Type() OpType {
	return TypeTransferToVesting
}

//Data returns the operation data TransferToVestingOperation.
func (op *TransferToVestingOperation) Data() interface{} {
	return op
}

//MarshalTransaction is a function of converting type TransferToVestingOperation to bytes.
func (op *TransferToVestingOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeTransferToVesting.Code()))
	enc.Encode(op.From)
	enc.Encode(op.To)
	enc.EncodeMoney(op.Amount)
	enc.EncodeMoney(op.Fee)
	return enc.Err()
}
