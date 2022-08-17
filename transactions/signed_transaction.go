//go:build !nosigning
// +build !nosigning

package transactions

import (
	"github.com/beowulf-foundation/beowulf-go/config"

	// Stdlib
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"time"

	// RPC
	"github.com/beowulf-foundation/beowulf-go/encoding/transaction"
	"github.com/beowulf-foundation/beowulf-go/types"

	// Vendor
	"github.com/pkg/errors"
)

//SignedTransaction structure of a signed transaction
type SignedTransaction struct {
	*types.Transaction
}

//NewSignedTransaction initialization of a new signed transaction
func NewSignedTransaction(tx *types.Transaction) *SignedTransaction {
	if tx.Expiration == nil {
		expiration := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
		tx.Expiration = &types.Time{Time: &expiration}
	}
	if tx.CreatedTime == 0 {
		createdTime := time.Now().UTC()
		tx.CreatedTime = types.UInt64(createdTime.Unix())
	}

	return &SignedTransaction{tx}
}

//Serialize function serializes a transaction
func (tx *SignedTransaction) Serialize() ([]byte, error) {
	var b bytes.Buffer
	encoder := transaction.NewEncoder(&b)

	if err := encoder.Encode(tx.Transaction); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//Digest function that returns a digest from a serialized transaction
func (tx *SignedTransaction) Digest(chain string) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode chain ID: %v", chain)
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Wrap(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Wrap(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	return digest[:], nil
}

//Sign function directly generating transaction signature, return transactionId
func (tx *SignedTransaction) Sign(privKeys [][]byte, chain string) (string, error) {
	var buf bytes.Buffer
	chainid, errdec := hex.DecodeString(chain)
	if errdec != nil {
		return "", errdec
	}

	txRaw, err := tx.Serialize()
	if err != nil {
		return "", err
	}
	// 	fmt.Println("Transaction:")
	// 	fmt.Println(tx)
	// 	fmt.Println("Hex of tx:")
	// 	tmp := hex.EncodeToString(txRaw)
	// 	fmt.Println(tmp)
	hashSha256 := sha256.Sum256(txRaw)
	var txId = make([]byte, 20)
	copy(txId, hashSha256[:20])
	buf.Write(chainid)
	buf.Write(txRaw)
	data := buf.Bytes()

	var sigsHex []string

	for _, privB := range privKeys {
		sigBytes, err := tx.SignSingle(privB, data)
		if err != nil {
			return "", err
		}
		sigsHex = append(sigsHex, hex.EncodeToString(sigBytes))
	}

	tx.Transaction.Signatures = sigsHex
	txHex := hex.EncodeToString(txId)
	return txHex, nil
}

func (tx *SignedTransaction) SignMulti(privKeys [][]byte, chain string) ([]string, error) {
	var sigsHex []string
	var buf bytes.Buffer
	chainid, errdec := hex.DecodeString(chain)
	if errdec != nil {
		return sigsHex, errdec
	}

	txRaw, err := tx.Serialize()
	if err != nil {
		return sigsHex, err
	}
	hashSha256 := sha256.Sum256(txRaw)
	var txId = make([]byte, 20)
	copy(txId, hashSha256[:20])
	buf.Write(chainid)
	buf.Write(txRaw)
	data := buf.Bytes()

	for _, privB := range privKeys {
		sigBytes, err := tx.SignSingle(privB, data)
		if err != nil {
			return sigsHex, err
		}
		sigsHex = append(sigsHex, hex.EncodeToString(sigBytes))
	}

	//tx.Transaction.Signatures = sigsHex
	return sigsHex, nil
}
