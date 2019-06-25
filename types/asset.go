package types

import (
	"encoding/json"
	"strconv"
	"strings"
	"beowulf-go/encoding/transaction"
)

//Asset type from parameter JSON
type Asset struct {
	Amount float64
	Symbol string
}

type AssetSymbol struct{
	Decimals uint8 		`json:"decimals"`
	AssetName string	`json:"name"`
}

//UnmarshalJSON unpacking the JSON parameter in the Asset type.
func (op *Asset) UnmarshalJSON(data []byte) error {
	str, errUnq := strconv.Unquote(string(data))
	if errUnq != nil {
		return errUnq
	}
	param := strings.Split(str, " ")

	s, errpf := strconv.ParseFloat(param[0], 64)
	if errpf != nil {
		return errpf
	}

	op.Amount = s
	op.Symbol = param[1]

	return nil
}

//MarshalJSON function for packing the Asset type in JSON.
func (op *Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(op.String())
}

//MarshalTransaction is a function of converting type Asset to bytes.
func (op *Asset) MarshalTransaction(encoder *transaction.Encoder) error {
	ans, err := json.Marshal(op)
	if err != nil {
		return err
	}

	str, err := strconv.Unquote(string(ans))
	if err != nil {
		return err
	}
	return encoder.EncodeMoney(str)
}

//String function convert type Asset to string.
func (op *Asset) String() string {
	var ammf string
	ammf = strconv.FormatFloat(op.Amount, 'f', 5, 64)
	return ammf + " " + op.Symbol
}

//StringAmount function convert type Asset.Amount to string.
func (op *Asset) StringAmount() string {
	return strconv.FormatFloat(op.Amount, 'f', 5, 64)
}

//UnmarshalJSON unpacking the JSON parameter in the AssetSymbol type.
func (op *AssetSymbol) UnmarshalJSON(data []byte) error {
	var raw AssetSymbol

	str := string (data)//strconv.Unquote(string(data))
	if str == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(str), &raw); err != nil {
		return err
	}

	op.Decimals = raw.Decimals
	op.AssetName = raw.AssetName
	return nil
}

//MarshalJSON function for packing the AssetSymbol type in JSON.
func (op *AssetSymbol) MarshalJSON() ([]byte, error) {
	ans, err := json.Marshal(*op)
	if err != nil {
		return []byte{}, err
	}
	return ans, nil
}

//MarshalTransaction is a function of converting type AssetSymbol to bytes.
func (op *AssetSymbol) MarshalTransaction(encoder *transaction.Encoder) error {
	ans, err := json.Marshal(*op)
	if err != nil {
		return err
	}

	str := string(ans)

	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeSymbol(str)
	return enc.Err()
}