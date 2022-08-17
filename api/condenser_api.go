package api

import (
	"encoding/json"

	"github.com/beowulf-foundation/beowulf-go/transports"
	"github.com/beowulf-foundation/beowulf-go/types"
	_ "github.com/beowulf-foundation/beowulf-go/types"
	"github.com/pkg/errors"
)

func (api *API) GetVersion() (*Version, error) {
	var resp Version
	err := api.call("condenser_api", "get_version", transports.EmptyParams, &resp, "")
	return &resp, err
}

//GetConfig api request get_config
func (api *API) GetConfig() (*Config, error) {
	var resp Config
	err := api.call("condenser_api", "get_config", transports.EmptyParams, &resp, "")
	return &resp, err
}

//GetDynamicGlobalProperties api request get_dynamic_global_properties
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call("condenser_api", "get_dynamic_global_properties", transports.EmptyParams, &resp, "")
	return &resp, err
}

//GetBlock api request get_block
func (api *API) GetBlock(blockNum uint32) (*Block, error) {
	var resp Block
	err := api.call("condenser_api", "get_block", []uint32{blockNum}, &resp, "")
	resp.Number = blockNum
	return &resp, err
}

//GetBlockHeader api request get_block_header
func (api *API) GetBlockHeader(blockNum uint32) (*BlockHeader, error) {
	var resp BlockHeader
	err := api.call("condenser_api", "get_block_header", []uint32{blockNum}, &resp, "")
	resp.Number = blockNum
	return &resp, err
}

// Set callback to invoke as soon as a new block is applied
func (api *API) SetBlockAppliedCallback(notice func(header *BlockHeader, error error)) (err error) {
	err = api.setCallback("condenser_api", "set_block_applied_callback", func(raw json.RawMessage) {
		var header []BlockHeader
		if err := json.Unmarshal(raw, &header); err != nil {
			notice(nil, err)
		}
		for _, b := range header {
			notice(&b, nil)
		}
	})
	return
}

func (api *API) GetSupernodeSchedule() (*SupernodeSchedule, error) {
	var resp SupernodeSchedule
	err := api.call("condenser_api", "get_supernode_schedule", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetHardforkVersion() (*string, error) {
	var resp string
	err := api.call("condenser_api", "get_hardfork_version", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetNextScheduledHardfork() (*ScheduledHardfork, error) {
	var resp ScheduledHardfork
	err := api.call("condenser_api", "get_next_scheduled_hardfork", transports.EmptyParams, &resp, "")
	return &resp, err
}

//GetTransaction api request get_transaction
func (api *API) GetTransaction(trxId string) (*TransactionResponse, error) {
	var resp TransactionResponse
	err := api.call("condenser_api", "get_transaction", []string{trxId}, &resp, "")
	//resp.ID = trxId
	return &resp, err
}

func (api *API) GetTransactionWithStatus(trxId string) (*TransactionResponse, error) {
	var resp TransactionResponse
	err := api.call("condenser_api", "get_transaction_with_status", []string{trxId}, &resp, "")
	//resp.ID = trxId
	return &resp, err
}

//GetTransactionHex api request get_transaction_hex
func (api *API) GetTransactionHex(tx *types.Transaction) (string, error) {
	var resp string
	err := api.call("condenser_api", "get_transaction_hex", []interface{}{tx}, &resp, "")
	return resp, err
}

//BroadcastTransaction api request broadcast_transaction
func (api *API) BroadcastTransaction(tx *types.Transaction) (*AsyncBroadcastResponse, error) {
	var resp AsyncBroadcastResponse
	err := api.call("condenser_api", "broadcast_transaction", []interface{}{tx}, &resp, "")
	return &resp, err
}

//BroadcastTransactionSynchronous api request broadcast_transaction_synchronous
func (api *API) BroadcastTransactionSynchronous(tx *types.Transaction) (*BroadcastResponse, error) {
	var resp BroadcastResponse
	err := api.call("condenser_api", "broadcast_transaction_synchronous", []interface{}{tx}, &resp, "")
	return &resp, err
}

func (api *API) GetAccounts(account string) (*AccountList, error) {
	var resp AccountList
	err := api.call("condenser_api", "get_accounts", [][]string{{account}}, &resp, "")
	return &resp, err
}

func (api *API) GetSupernodes(id types.UInt16) (*SupernodeList, error) {
	var resp SupernodeList
	err := api.call("condenser_api", "get_supernodes", [][]types.UInt16{{id}}, &resp, "")
	return &resp, err
}

func (api *API) GetSupernodeByAccount(account string) (*SupernodeInfo, error) {
	var resp SupernodeInfo
	err := api.call("condenser_api", "get_supernode_by_account", []string{account}, &resp, "")
	return &resp, err
}

func (api *API) GetSupernodeByVote(lowerBound string, limit uint32) (*SupernodeList, error) {
	var resp SupernodeList
	err := api.call("condenser_api", "get_supernodes_by_vote", []interface{}{lowerBound, limit}, &resp, "")
	return &resp, err
}

func (api *API) LookupSupernodeAccounts(lowerBound string, limit uint32) (*[]string, error) {
	var resp []string
	err := api.call("condenser_api", "lookup_supernode_accounts", []interface{}{lowerBound, limit}, &resp, "")
	return &resp, err
}

func (api *API) GetSupernodeCount() (*uint64, error) {
	var resp uint64
	err := api.call("condenser_api", "get_supernode_count", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetSupernodeVoted(account string) (*SupernodeVoteList, error) {
	var resp SupernodeVoteList
	err := api.call("condenser_api", "get_supernode_voted_by_acc", []string{account}, &resp, "")
	return &resp, err
}

func (api *API) GetKeyReferences(publicKey string) (*[][]string, error) {
	var resp [][]string
	err := api.call("condenser_api", "get_key_references", [][]string{{publicKey}}, &resp, "")
	return &resp, err
}

func (api *API) ListAccounts(lowerBound string, limit uint32) (*[]string, error) {
	var resp []string
	err := api.call("condenser_api", "lookup_accounts", []interface{}{lowerBound, limit}, &resp, "")
	return &resp, err
}

func (api *API) GetAccountCount() (*uint64, error) {
	var resp uint64
	err := api.call("condenser_api", "get_account_count", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetActiveSupernodes() (*[]string, error) {
	var resp []string
	err := api.call("condenser_api", "get_active_supernodes", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) ListSupernodes(lowerBound string, limit uint32) (*[]string, error) {
	var resp []string
	err := api.call("condenser_api", "lookup_supernode_accounts", []interface{}{lowerBound, limit}, &resp, "")
	return &resp, err
}

func (api *API) ListTokens() (*TokenList, error) {
	var resp TokenList
	err := api.call("condenser_api", "list_smt_tokens", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetTokens(name string) (*TokenList, error) {
	var resp TokenList
	err := api.call("condenser_api", "find_smt_tokens_by_name", []string{name}, &resp, "")
	return &resp, err
}

//Get balance
func (api *API) GetBalance(account, tokenName string, decimals uint8) (*string, error) {
	var resp string
	var symbol types.AssetSymbol
	symbol.AssetName = tokenName
	symbol.Decimals = decimals
	err := api.call("condenser_api", "get_balance", []interface{}{account, symbol}, &resp, "")
	return &resp, err
}

func (api *API) GetPendingTransactionCount() (*uint64, error) {
	var resp uint64
	err := api.call("condenser_api", "get_pending_transaction_count", transports.EmptyParams, &resp, "")
	return &resp, err
}

func (api *API) GetNFTs(symbol string, limit, offset uint32) (*NFTList, error) {
	var resp NFTList
	var params Params
	params.Contract = "nft"
	params.Table = "nfts"
	obj := map[string]interface{}{}
	if len(symbol) > 0 {
		obj["symbol"] = symbol
	}
	params.Query = obj
	params.Limit = limit
	params.Offset = offset
	err := api.call("", "find", params, &resp, "s01")
	return &resp, err
}

func (api *API) GetNFTBalance(account string, symbol string, limit, offset uint32) (*NFTInstanceList, error) {
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol can't be null")
	}
	var resp NFTInstanceList
	var params Params
	params.Contract = "nft"
	params.Table = symbol + "instances"
	obj := map[string]interface{}{}
	if len(account) > 0 {
		obj["account"] = account
	}
	params.Query = obj
	params.Limit = limit
	params.Offset = offset
	err := api.call("", "find", params, &resp, "s01")
	return &resp, err
}

func (api *API) GetNFTInstances(symbol string, limit, offset uint32) (*NFTInstanceList, error) {
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol can't be null")
	}
	var resp NFTInstanceList
	var params Params
	params.Contract = "nft"
	params.Table = symbol + "instances"
	obj := map[string]interface{}{}
	params.Query = obj
	params.Limit = limit
	params.Offset = offset
	err := api.call("", "find", params, &resp, "s01")
	return &resp, err
}

func (api *API) GetNFTBalanceOfAccount(account string, limit, offset uint32) (map[string]NFTInstanceList, error) {
	result := make(map[string]NFTInstanceList)
	var res NFTList
	var params Params
	params.Contract = "nft"
	params.Table = "nfts"
	obj := map[string]interface{}{}
	params.Query = obj
	err := api.call("", "find", params, &res, "s01")
	for _, element := range res {
		var resp NFTInstanceList
		symbol := element.Symbol
		var instance NFTInstanceList
		var params Params
		params.Contract = "nft"
		params.Table = symbol + "instances"
		obj := map[string]interface{}{}
		if len(account) > 0 {
			obj["account"] = account
		}
		params.Query = obj
		params.Limit = limit
		params.Offset = offset
		err = api.call("", "find", params, &instance, "s01")
		if err == nil && len(instance) > 0 {
			resp = append(resp, instance...)
		}
		result[symbol] = resp
	}
	return result, err
}
func (api *API) GetLatestNFTBlock() (*NFTBlock, error) {
	var resp NFTBlock
	err := api.call("", "getLatestBlockInfo", transports.EmptyParams, &resp, "s01")
	//resp.Number = blockNum
	return &resp, err
}

func (api *API) GetNFTBlock(blockNum uint32) (*NFTBlock, error) {
	var resp NFTBlock
	var params BlockParams
	params.BlockNumber = blockNum
	err := api.call("", "getBlockInfo", params, &resp, "s01")
	//resp.Number = blockNum
	return &resp, err
}

func (api *API) GetNFTTransaction(trxId string) (*NFTTransaction, error) {
	var resp NFTTransaction
	var params TransactionParams
	params.Txid = trxId
	err := api.call("", "getTransactionInfo", params, &resp, "s01")
	//resp.ID = trxId
	return &resp, err
}
