package types

// OpType represents a Golos operation type, i.e. vote, comment, pow and so on.
type OpType string

// Code returns the operation code associated with the given operation type.
func (kind OpType) Code() uint16 {
	return opCodes[kind]
}

const (
	TypeTransfer                    OpType = "transfer"
	TypeTransferToVesting           OpType = "transfer_to_vesting"
	TypeWithdrawVesting             OpType = "withdraw_vesting"
	TypeAccountCreate               OpType = "account_create"
	TypeAccountUpdate               OpType = "account_update"
	TypeSupernodeUpdate             OpType = "supernode_update"
	TypeAccountSupernodeVote        OpType = "account_supernode_vote"
	TypeSmtCreate					OpType = "smt_create"
	TypeFillVestingWithdraw         OpType = "fill_vesting_withdraw"      //Virtual Operation
	TypeShutdownSupernode           OpType = "shutdown_supernode"         //Virtual Operation
	TypeHardfork                    OpType = "hardfork"                   //Virtual Operation
	TypeProducerReward    	        OpType = "producer_reward"  		  //Virtual Operation
	TypeClearNullAccountBalance     OpType = "clear_null_account_balance" //Virtual Operation
)

var opTypes = [...]OpType{
	TypeTransfer,
	TypeTransferToVesting,
	TypeWithdrawVesting,
	TypeAccountCreate,
	TypeAccountUpdate,
	TypeSupernodeUpdate,
	TypeAccountSupernodeVote,
	TypeSmtCreate,
	TypeFillVestingWithdraw,     //Virtual Operation
	TypeShutdownSupernode,       //Virtual Operation
	TypeHardfork,                //Virtual Operation
	TypeProducerReward,     	 //Virtual Operation
	TypeClearNullAccountBalance, //Virtual Operation
}

// opCodes keeps mapping operation type -> operation code.
var opCodes map[OpType]uint16

func init() {
	opCodes = make(map[OpType]uint16, len(opTypes))
	for i, opType := range opTypes {
		opCodes[opType] = uint16(i)
	}
}
