package types

import (
	"errors"
	"beowulf-go/encoding/transaction"
)

// Transaction represents a blockchain transaction.
type Transaction struct {
	//transaction_id_type        transaction_id;
	//uint32_t                   block_num = 0;
	//uint32_t                   transaction_num = 0;
	RefBlockNum    UInt16     `json:"ref_block_num"`
	RefBlockPrefix UInt32     `json:"ref_block_prefix"`
	Expiration     *Time      `json:"expiration"`
	Operations     Operations `json:"operations"`
	Extensions     [][]interface{}      `json:"extensions"`
	CreatedTime	   UInt64	  `json:"created_time"`
	Signatures     []string   `json:"signatures"`
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
	enc.EncodeUVarint(0)
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
