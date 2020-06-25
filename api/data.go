package api

import (
	"beowulf-go/types"
)

//BroadcastResponse structure for the BroadcastTransactionSynchronous function
type BroadcastResponse struct {
	ID          string `json:"id"`
	BlockNum    int32  `json:"block_num"`
	TrxNum      int32  `json:"trx_num"`
	Expired     bool   `json:"expired"`
	CreatedTime int64  `json:"created_time"`
}

type AsyncBroadcastResponse struct {
	ID string `json:"id"`
}

//Version structure for the GetVersion function.
type Version struct {
	BlockchainVersion string `json:"blockchain_version"`
	BeowulfRevision   string `json:"beowulf_revision"`
	FcRevision        string `json:"fc_revision"`
}

//Config structure for the GetConfig function.
type Config struct {
	TokenCreationFee int `json:"SMT_TOKEN_CREATION_FEE"`
	//WDSymbol						 string		   `json:"WD_SYMBOL"`
	Percent100                            *types.Int  `json:"BEOWULF_100_PERCENT"`
	Percent1                              *types.Int  `json:"BEOWULF_1_PERCENT"`
	AddressPrefix                         string      `json:"BEOWULF_ADDRESS_PREFIX"`
	HardforkVersion                       string      `json:"BEOWULF_BLOCKCHAIN_HARDFORK_VERSION"`
	Version                               string      `json:"BEOWULF_BLOCKCHAIN_VERSION"`
	BlockInterval                         uint        `json:"BEOWULF_BLOCK_INTERVAL"`
	BlocksPerDay                          *types.Int  `json:"BEOWULF_BLOCKS_PER_DAY"`
	BlocksPerYear                         *types.Int  `json:"BEOWULF_BLOCKS_PER_YEAR"`
	ChainID                               string      `json:"BEOWULF_CHAIN_ID"`
	GenesisTime                           *types.Time `json:"BEOWULF_GENESIS_TIME"`
	HardforkRequiredSupernodes            *types.Int  `json:"BEOWULF_HARDFORK_REQUIRED_SUPERNODES"`
	InflationNarrowingPeriod              *types.Int  `json:"BEOWULF_INFLATION_NARROWING_PERIOD"`
	InflationRateStartPercent             *types.Int  `json:"BEOWULF_INFLATION_RATE_START_PERCENT"`
	InflationRateStopPercent              *types.Int  `json:"BEOWULF_INFLATION_RATE_STOP_PERCENT"`
	InitiatorName                         string      `json:"BEOWULF_INIT_MINER_NAME"`
	InitiatorPublicKey                    string      `json:"BEOWULF_INIT_PUBLIC_KEY_STR"`
	InitSupply                            *types.Int  `json:"BEOWULF_INIT_SUPPLY"`
	WDInitSupply                          *types.Int  `json:"WD_INIT_SUPPLY"`
	IrreversibleThreshold                 *types.Int  `json:"BEOWULF_IRREVERSIBLE_THRESHOLD"`
	MaxAccountNameLength                  *types.Int  `json:"BEOWULF_MAX_ACCOUNT_NAME_LENGTH"`
	MaxAccountSupernodeVotes              *types.Int  `json:"BEOWULF_MAX_ACCOUNT_SUPERNODE_VOTES"`
	MaxAuthorityMembership                *types.Int  `json:"BEOWULF_MAX_AUTHORITY_MEMBERSHIP"`
	BlockSize                             *types.Int  `json:"BEOWULF_SOFT_MAX_BLOCK_SIZE"`
	MaxMemoSize                           *types.Int  `json:"BEOWULF_MAX_MEMO_SIZE"`
	MaxSupernodes                         *types.Int  `json:"BEOWULF_MAX_SUPERNODES"`
	MaxPermanentSupernodes                *types.Int  `json:"BEOWULF_MAX_PERMANENT_SUPERNODES_HF0"`
	MaxRunnerSupernodes                   *types.Int  `json:"BEOWULF_MAX_RUNNER_SUPERNODES_HF0"`
	MaxShareSupply                        *types.Int  `json:"BEOWULF_MAX_SHARE_SUPPLY"`
	MaxSigCheckDepth                      *types.Int  `json:"BEOWULF_MAX_SIG_CHECK_DEPTH"`
	MaxSigCheckAccounts                   *types.Int  `json:"BEOWULF_MAX_SIG_CHECK_ACCOUNTS"`
	MaxTimeUntilExpiration                int         `json:"BEOWULF_MAX_TIME_UNTIL_EXPIRATION"`
	MaxTransactionSize                    *types.Int  `json:"BEOWULF_MAX_TRANSACTION_SIZE"`
	MaxUndoHistory                        *types.Int  `json:"BEOWULF_MAX_UNDO_HISTORY"`
	MaxVotedSupernodes                    *types.Int  `json:"BEOWULF_MAX_VOTED_SUPERNODES_HF0"`
	MinSupernodeFund                      *types.Int  `json:"BEOWULF_MIN_SUPERNODE_FUND"`
	MinTransactionFee                     *types.Int  `json:"BEOWULF_MIN_TRANSACTION_FEE"`
	MinAccountCreationFee                 *types.Int  `json:"BEOWULF_MIN_ACCOUNT_CREATION_FEE"`
	MinAccountNameLength                  *types.Int  `json:"BEOWULF_MIN_ACCOUNT_NAME_LENGTH"`
	MinBlockSizeLimit                     *types.Int  `json:"BEOWULF_MIN_BLOCK_SIZE"`
	NullAccount                           string      `json:"BEOWULF_NULL_ACCOUNT"`
	NumInitiators                         *types.Int  `json:"BEOWULF_NUM_INIT_MINERS"`
	OwnerAuthHistoryTrackingStartBlockNum *types.Int  `json:"BEOWULF_OWNER_AUTH_HISTORY_TRACKING_START_BLOCK_NUM"`
	OwnerUpdateLimit                      *types.Int  `json:"BEOWULF_OWNER_UPDATE_LIMIT"`
	VestingWithdrawIntervals              *types.Int  `json:"BEOWULF_VESTING_WITHDRAW_INTERVALS"`
	VestingWithdrawIntervalSeconds        *types.Int  `json:"BEOWULF_VESTING_WITHDRAW_INTERVAL_SECONDS"`
	//BEOWULF_SYMBOL
	//VESTS_SYMBOL
	VirtualScheduleLapLength  *types.Int `json:"BEOWULF_VIRTUAL_SCHEDULE_LAP_LENGTH2"`
	Beowulf1Beowulf           *types.Int `json:"BEOWULF_1_BEOWULF"`
	Beowulf1Vests             *types.Int `json:"BEOWULF_1_VESTS"`
	MaxTokenPerAccount        *types.Int `json:"BEOWULF_MAX_TOKEN_PER_ACCOUNT"`
	MinPermanentSupernodeFund *types.Int `json:"BEOWULF_MIN_PERMANENT_SUPERNODE_FUND"`
	MaxTokenNameLength        *types.Int `json:"BEOWULF_MAX_TOKEN_NAME_LENGTH"`
	MinTokenNameLength        *types.Int `json:"BEOWULF_MIN_TOKEN_NAME_LENGTH"`
	SymbolBeowulf             string     `json:"BEOWULF_SYMBOL_BEOWULF"`
	SymbolWD                  string     `json:"BEOWULF_SYMBOL_WD"`
	SymbolVests               string     `json:"BEOWULF_SYMBOL_VESTS"`
	BlockRewardGap            *types.Int `json:"BEOWULF_BLOCK_REWARD_GAP"`
}

//DynamicGlobalProperties structure for the GetDynamicGlobalProperties function.
type DynamicGlobalProperties struct {
	ID               string       `json:"id"`
	HeadBlockNumber  uint32       `json:"head_block_number"`
	HeadBlockID      string       `json:"head_block_id"`
	Time             *types.Time  `json:"time"`
	CurrentSupernode string       `json:"current_supernode"`
	CurrentSupply    *types.Asset `json:"current_supply"`
	CurrentWDSupply  *types.Asset `json:"current_wd_supply"`

	TotalVestingFund         *types.Asset `json:"total_vesting_fund_beowulf"`
	TotalVestingShares       *types.Asset `json:"total_vesting_shares"`
	CurrentAslot             uint64       `json:"current_aslot"`
	//RecentSlotsFilled        *types.Int   `json:"recent_slots_filled"`
	//ParticipationCount       uint8        `json:"participation_count"`
	LastIrreversibleBlockNum uint32       `json:"last_irreversible_block_num"`
}

//BlockHeader structure for the GetBlockHeader and SetBlockAppliedCallback functions
type BlockHeader struct {
	Number                uint32        `json:"-"`
	Previous              string        `json:"previous"`
	Timestamp             string        `json:"timestamp"`
	Supernode             string        `json:"supernode"`
	TransactionMerkleRoot string        `json:"transaction_merkle_root"`
	BlockReward           *types.Asset  `json:"block_reward"`
	Extensions            []interface{} `json:"extensions"`
}

//Block structure for the GetBlock function
type Block struct {
	Number                uint32               `json:"-"`
	Previous              string               `json:"previous"`
	Timestamp             *types.Time          `json:"timestamp"`
	Supernode             string               `json:"supernode"`
	TransactionMerkleRoot string               `json:"transaction_merkle_root"`
	BlockReward           *types.Asset         `json:"block_reward"`
	Extensions            [][]interface{}      `json:"extensions"`
	SupernodeSignature    string               `json:"supernode_signature"`
	Transactions          []*types.Transaction `json:"transactions"`
	BlockId               string               `json:"block_id"`
	PublicKeyType         string               `json:"signing_key"`
	TransactionIds        []string             `json:"transaction_ids"`
}

//SupernodeSchedule structure for the GetSupernodeSchedule function
type SupernodeSchedule struct {
	Id                              *types.UInt16 `json:"id"`
	CurrentVirtualTime              string        `json:"current_virtual_time"`
	NextShuffleBlockNum             uint32        `json:"next_shuffle_block_num"`
	CurrentShuffledSupernodes       []string      `json:"current_shuffled_supernodes"`
	NumScheduledSupernodes          uint8         `json:"num_scheduled_supernodes"`
	ElectedWeight                   uint8         `json:"elected_weight"`
	TimeshareWeight                 uint8         `json:"timeshare_weight"`
	PermanentWeight                 uint8         `json:"permanent_weight"`
	SupernodePayNormalizationFactor uint32        `json:"supernode_pay_normalization_factor"`
	MajorityVersion                 string        `json:"majority_version"`
	MaxVotedSupernodes              uint8         `json:"max_voted_supernodes"`
	MaxPermanentSupernodes          uint8         `json:"max_permanent_supernodes"`
	MaxRunnerSupernodes             uint8         `json:"max_runner_supernodes"`
	HardforkRequiredSupernodes      uint8         `json:"hardfork_required_supernodes"`
}

type ScheduledHardfork struct {
	HardforkVersion string      `json:"hf_version"`
	LiveTime        *types.Time `json:"live_time"`
}

type AccountInfo struct {
	Id                    *types.UInt16          `json:"id"`
	Name                  string                 `json:"name"`
	Owner                 *types.Authority       `json:"owner"`
	JSONMetadata          *types.AccountMetadata `json:"json_metadata"`
	LastOwnerUpdate       *types.Time            `json:"last_owner_update"`
	LastAccountUpdate     *types.Time            `json:"last_account_update"`
	Created               *types.Time            `json:"created"`
	Balance               string                 `json:"balance"`
	WdBalance             string                 `json:"wd_balance"`
	VestingShares         string                 `json:"vesting_shares"`
	VestingWithdrawRate   string                 `json:"vesting_withdraw_rate"`
	NextVestingWithdrawal *types.Time            `json:"next_vesting_withdrawal"`
	Withdrawn             *types.Int64           `json:"withdrawn"`
	ToWithdraw            *types.Int64           `json:"to_withdraw"`
	SupernodesVotedFor    uint16                 `json:"supernodes_voted_for"`
	TokenList             []string               `json:"token_list"`
	VestingBalance        string                 `json:"vesting_balance"`
	SupernodeVotes        []string               `json:"supernode_votes"`
}

type AccountList []AccountInfo

type SupernodeInfo struct {
	Id                    *types.UInt16 `json:"id"`
	Owner                 string        `json:"owner"`
	Created               *types.Time   `json:"created"`
	TotalMissed           uint32        `json:"total_missed"`
	LastASlot             uint64        `json:"last_aslot"`
	LastConfirmedBlockNum uint64        `json:"last_confirmed_block_num"`
	SigningKey            string        `json:"signing_key"`
	Votes                 *types.Int64  `json:"votes"`
	//VirtualLastUpdate 		*big.Int			`json:"virtual_last_update"`
	//VirtualPosition			*big.Int			`json:"virtual_position"`
	//VirtualScheduledTime	*big.Int			`json:"virtual_scheduled_time"`
	LastWork            string      `json:"last_work"`
	RunningVersion      string      `json:"running_version"`
	HardforkVersionVote string      `json:"hardfork_version_vote"`
	HardforkTimeVote    *types.Time `json:"hardfork_time_vote"`
}

type SupernodeList []SupernodeInfo

type SupernodeVoteInfo struct {
	Id        *types.UInt16 `json:"id"`
	Supernode string        `json:"supernode"`
	Account   string        `json:"account"`
	Votes     *types.Int64  `json:"votes"`
}

type SupernodeVoteList []SupernodeVoteInfo

type TokenInfo struct {
	Id             *types.UInt16      `json:"id"`
	LiquidSymbol   *types.AssetSymbol `json:"liquid_symbol"`
	ControlAccount string             `json:"control_account"`
	Phase          uint8              `json:"phase"`
	CurrentSupply  *types.Int64       `json:"current_supply"`
}

type TokenList []TokenInfo

type TransactionResponse struct {
	RefBlockNum    *types.UInt16     `json:"ref_block_num"`
	RefBlockPrefix *types.UInt32     `json:"ref_block_prefix"`
	Expiration     *types.Time       `json:"expiration"`
	Operations     *types.Operations `json:"operations"`
	Extensions     []interface{}     `json:"extensions"`
	CreatedTime    *types.UInt64     `json:"created_time"`
	Signatures     []string          `json:"signatures"`
	TransactionId  string            `json:"transaction_id"`
	BlockNum       *types.UInt32     `json:"block_num"`
	TransactionNum *types.UInt32     `json:"transaction_num"`
	Status         string            `json:"status"`
}

