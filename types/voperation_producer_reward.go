package types

//LiquidityRewardOperation represents liquidity_reward operation data.
type ProducerRewardOperation struct {
	Producer  		string `json:"producer"`
	VestingShares 	*Asset `json:"vesting_shares"`
}

//Type function that defines the type of operation LiquidityRewardOperation.
func (op *ProducerRewardOperation) Type() OpType {
	return TypeProducerReward
}

//Data returns the operation data LiquidityRewardOperation.
func (op *ProducerRewardOperation) Data() interface{} {
	return op
}
