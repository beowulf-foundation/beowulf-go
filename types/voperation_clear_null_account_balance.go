package types

//ReturnVestingDelegationOperation represents return_vesting_delegation operation data.
type ClearNullAccountBalanceOperation struct {
	TotalCleared       []*Asset `json:"total_cleared"`
}

//Type function that defines the type of operation ReturnVestingDelegationOperation.
func (op *ClearNullAccountBalanceOperation) Type() OpType {
	return TypeClearNullAccountBalance
}

//Data returns the operation data ReturnVestingDelegationOperation.
func (op *ClearNullAccountBalanceOperation) Data() interface{} {
	return op
}
