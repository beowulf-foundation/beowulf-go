package types

import (
	"beowulf-go/encoding/transaction"
)

//Authority is an additional structure used by other structures.
type Authority struct {
	AccountAuths    StringInt64Map `json:"account_auths"`
	KeyAuths        StringInt64Map `json:"key_auths"`
	WeightThreshold uint32         `json:"weight_threshold"`
}

//MarshalTransaction is a function of converting type Authority to bytes.
func (auth *Authority) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeNumber(uint32(auth.WeightThreshold))
	// encode AccountAuths as map[string]uint16
	enc.EncodeUVarint(uint64(len(auth.AccountAuths)))
	acc_keys := []string{}
	for acc_key := range auth.AccountAuths {
		acc_keys = append(acc_keys, acc_key)
	}
	//acc_tmp := make(map[string]int64)
	//for i:=len(acc_keys)-1; i>=0; i--{
	//	item := acc_keys[i]
	//	acc_tmp[item] = auth.AccountAuths[item]
	//}
	for k, v := range auth.AccountAuths {
		enc.EncodeString(k)
		enc.EncodeNumber(uint16(v))
	}
	// encode KeyAuths as map[PubKey]uint16
	enc.EncodeUVarint(uint64(len(auth.KeyAuths)))
	keys := []string{}
	for key := range auth.KeyAuths {
		keys = append(keys, key)
	}
	//tmp := make(map[string]int64)
	//for i:=len(keys)-1; i>=0; i--{
	//	item := keys[i]
	//	tmp[item] = auth.KeyAuths[item]
	//}
	for k, v := range auth.KeyAuths {
		enc.EncodePubKey(k)
		enc.EncodeNumber(uint16(v))
	}
	return enc.Err()
}
