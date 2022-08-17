package types

import (
	"errors"

	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	RefBlockNum    UInt16        `json:"ref_block_num"`
	RefBlockPrefix UInt32        `json:"ref_block_prefix"`
	Expiration     *Time         `json:"expiration"`
	Operations     Operations    `json:"operations"`
	Extensions     []interface{} `json:"extensions"`
	CreatedTime    UInt64        `json:"created_time"`
	Signatures     []string      `json:"signatures"`
}

// MarshalTransaction implements transaction.Marshaller interface.
func (tx *Transaction) MarshalTransaction(encoder *transaction.Encoder) error {
	if len(tx.Operations) == 0 {
		return errors.New("no operation specified")
	}

	enc := transaction.NewRollingEncoder(encoder)

	enc.Encode(tx.RefBlockNum)
	enc.Encode(tx.RefBlockPrefix)
	enc.Encode(tx.Expiration)

	enc.EncodeUVarint(uint64(len(tx.Operations)))
	for _, op := range tx.Operations {
		enc.Encode(op)
	}

	// Extensions are not supported yet.
	if len(tx.Extensions) > 0 {
		enc.EncodeUVarint(uint64(len(tx.Extensions)))
		for _, ext := range tx.Extensions {
			enc.Encode(ext)
		}
	} else {
		enc.EncodeUVarint(0)
	}

	enc.Encode(tx.CreatedTime)
	for _, sig := range tx.Signatures {
		enc.Encode(sig)
	}
	return enc.Err()
}

// PushOperation can be used to add an operation into the transaction.
func (tx *Transaction) PushOperation(op Operation) {
	tx.Operations = append(tx.Operations, op)
}
