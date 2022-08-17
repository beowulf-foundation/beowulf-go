package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beowulf-foundation/beowulf-go/api"
	"github.com/beowulf-foundation/beowulf-go/config"
	"github.com/beowulf-foundation/beowulf-go/transactions"
	"github.com/beowulf-foundation/beowulf-go/types"
)

func (client *Client) GetBlock(blockNum uint32) (*api.Block, error) {
	return client.API.GetBlock(blockNum)
}

func (client *Client) GetTransaction(trx string) (*api.TransactionResponse, error) {
	return client.API.GetTransaction(trx)
}

func (client *Client) GetAccount(account string) (*api.AccountInfo, error) {
	accounts, err := client.API.GetAccounts(account)
	if err != nil {
		return nil, err
	}
	var list []api.AccountInfo
	list = *accounts
	if len(list) == 0 {
		return nil, errors.New("Unknown account")
	}
	return &list[0], nil
}

func (client *Client) GetSupernodeByAccount(account string) (*api.SupernodeInfo, error) {
	return client.API.GetSupernodeByAccount(account)
}

func (client *Client) GetSupernodeVoted(account string) (*api.SupernodeVoteList, error) {
	return client.API.GetSupernodeVoted(account)
}

func (client *Client) GetKeyReferences(publicKey string) (*[][]string, error) {
	return client.API.GetKeyReferences(publicKey)
}

func (client *Client) ListAccounts(lowerBound string, limit uint32) (*[]string, error) {
	return client.API.ListAccounts(lowerBound, limit)
}

func (client *Client) GetActiveSupernodes() (*[]string, error) {
	return client.API.GetActiveSupernodes()
}

func (client *Client) ListSupernodes(lowerBound string, limit uint32) (*[]string, error) {
	return client.API.ListSupernodes(lowerBound, limit)
}

func (client *Client) ListTokens() (*api.TokenList, error) {
	return client.API.ListTokens()
}

func (client *Client) GetToken(name string) (*api.TokenInfo, error) {
	tokens, err := client.API.GetTokens(name)
	if err != nil {
		return nil, err
	}
	var list []api.TokenInfo
	list = *tokens
	if len(list) == 0 {
		return nil, errors.New("Unknown token")
	}
	return &list[0], nil
}

func (client *Client) GetBalance(account, tokenName string, decimals uint8) (*string, error) {
	return client.API.GetBalance(account, tokenName, decimals)
}

func (client *Client) CommitScb(scid, committer, content, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if !validate {
		return nil, errors.New("Fee is not valid")
	}

	var trx []types.Operation
	tx := &types.ScbValidateOperation{
		Committer:   committer,
		Scid:        scid,
		ScOperation: content,
		Fee:         fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "ScbValidate", Bresp: resp}, err
}

func (client *Client) GetNFTs(symbol string, limit, offset uint32) (*api.NFTList, error) {
	return client.API.GetNFTs(symbol, limit, offset)
}

func (client *Client) GetNFTBalance(account, symbol string, limit, offset uint32) (*api.NFTInstanceList, error) {
	return client.API.GetNFTBalance(account, symbol, limit, offset)
}

func (client *Client) GetNFTInstances(symbol string, limit, offset uint32) (*api.NFTInstanceList, error) {
	return client.API.GetNFTInstances(symbol, limit, offset)
}

func (client *Client) GetNFTBalanceOfAccount(account string, limit, offset uint32) (map[string]api.NFTInstanceList, error) {
	return client.API.GetNFTBalanceOfAccount(account, limit, offset)
}

func (client *Client) GetLatestNFTBlock() (*api.NFTBlock, error) {
	return client.API.GetLatestNFTBlock()
}

func (client *Client) GetNFTBlock(blockNum uint32) (*api.NFTBlock, error) {
	return client.API.GetNFTBlock(blockNum)
}

func (client *Client) GetNFTTransaction(trxid string) (*api.NFTTransaction, error) {
	return client.API.GetNFTTransaction(trxid)
}

//Create NFT
func (client *Client) CreateNFT(fromName, scid, name, symbol, maxSupply, fee string, authorizedIssuingAccounts []string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(name) <= 0 {
		return nil, errors.New("Name is not valid")
	}
	validateSymbol := ValidateNftSymbol(symbol)
	if validateSymbol != nil {
		return nil, validateSymbol
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	accounts := ""
	for i, item := range authorizedIssuingAccounts {
		if i == 0 {
			accounts += "[\"" + item + "\""
			if i == len(authorizedIssuingAccounts)-1 {
				accounts += "]"
			}
		} else if i == len(authorizedIssuingAccounts)-1 {
			accounts += "," + "\"" + item + "\"]"
		} else {
			accounts += "," + "\"" + item + "\""
		}
	}

	var scoperation string
	if len(authorizedIssuingAccounts) > 0 {
		if len(maxSupply) > 0 {
			scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"create\",\"contractPayload\":{\"name\":\"%s\",\"symbol\":\"%s\",\"maxSupply\":\"%s\",\"authorizedIssuingAccounts\":%s}}", name, symbol, maxSupply, accounts)
		} else {
			scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"create\",\"contractPayload\":{\"name\":\"%s\",\"symbol\":\"%s\",\"authorizedIssuingAccounts\":%s}}", name, symbol, accounts)
		}
	} else {
		if len(maxSupply) > 0 {
			scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"create\",\"contractPayload\":{\"name\":\"%s\",\"symbol\":\"%s\",\"maxSupply\":\"%s\"}}", name, symbol, maxSupply)
		} else {
			scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"create\",\"contractPayload\":{\"name\":\"%s\",\"symbol\":\"%s\"}}", name, symbol)
		}
	}
	//scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"create\",\"contractPayload\":{\"name\":\"%s\",\"symbol\":\"%s\",\"maxSupply\":\"%s\"}}", name, symbol, maxSupply)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

//func (client *Client) UpdateUrl(fromName, scid, symbol, url, fee string) (*OperResp, error) {
//	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
//	if validate == false {
//		return nil, errors.New("Fee is not valid")
//	}
//	var owners []string
//	owners = append(owners, fromName)
//
//	var scoperation string
//	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updateMetadata\",\"contractPayload\":{\"symbol\":\"%s\",\"metadata\":{\"url\":\"%s\"}}}", symbol, url)
//	var trx []types.Operation
//	tx := &types.SmartContractOperation{
//		RequiredOwners: owners,
//		Scid:           scid,
//		ScOperation:    scoperation,
//		Fee:            fee,
//	}
//	trx = append(trx, tx)
//	resp, err := client.SendTrx(trx, "")
//	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
//}
//
//func (client *Client) UpdateImage(fromName, scid, symbol, image, fee string) (*OperResp, error) {
//	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
//	if validate == false {
//		return nil, errors.New("Fee is not valid")
//	}
//	var owners []string
//	owners = append(owners, fromName)
//
//	var scoperation string
//	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updateMetadata\",\"contractPayload\":{\"symbol\":\"%s\",\"metadata\":{\"image\":\"%s\"}}}", symbol, image)
//	var trx []types.Operation
//	tx := &types.SmartContractOperation{
//		RequiredOwners: owners,
//		Scid:           scid,
//		ScOperation:    scoperation,
//		Fee:            fee,
//	}
//	trx = append(trx, tx)
//	resp, err := client.SendTrx(trx, "")
//	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
//}

func (client *Client) UpdateMetadata(fromName, scid, symbol, url, image, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updateMetadata\",\"contractPayload\":{\"symbol\":\"%s\",\"metadata\":{\"url\":\"%s\",\"image\":\"%s\"}}}", symbol, url, image)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) UpdateName(fromName, scid, symbol, name, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updateName\",\"contractPayload\":{\"symbol\":\"%s\",\"name\":\"%s\"}}", symbol, name)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) UpdateOrgName(fromName, scid, symbol, orgName, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updateOrgName\",\"contractPayload\":{\"symbol\":\"%s\",\"orgName\":\"%s\"}}", symbol, orgName)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) AddProperty(fromName, scid, symbol, propertyName, propertyType, fee string, authorizedEditingAccounts []string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(propertyName) <= 0 || len(propertyType) <= 0 {
		return nil, errors.New("Property is not valid")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	accounts := ""
	for i, item := range authorizedEditingAccounts {
		if i == 0 {
			accounts += "[\"" + item + "\""
			if i == len(authorizedEditingAccounts)-1 {
				accounts += "]"
			}
		} else if i == len(authorizedEditingAccounts)-1 {
			accounts += "," + "\"" + item + "\"]"
		} else {
			accounts += "," + "\"" + item + "\""
		}
	}
	var scoperation string
	if len(authorizedEditingAccounts) > 0 {
		scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"addProperty\",\"contractPayload\":{\"symbol\":\"%s\",\"name\":\"%s\",\"type\":\"%s\",\"authorizedEditingAccounts\":%s}}", symbol, propertyName, propertyType, accounts)
	} else {
		scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"addProperty\",\"contractPayload\":{\"symbol\":\"%s\",\"name\":\"%s\",\"type\":\"%s\"}}", symbol, propertyName, propertyType)
	}
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) IssueNFT(fromName, scid, symbol, to, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(to) <= 0 {
		return nil, errors.New("Recipient is not valid")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}

	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"issue\",\"contractPayload\":{\"symbol\":\"%s\",\"to\":\"%s\",\"toType\":\"user\",\"feeSymbol\":\"BEE\"}}", symbol, to)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) IssueWithProperties(fromName, scid, symbol, to, fee string, properties interface{}) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var owners []string
	owners = append(owners, fromName)
	if len(to) <= 0 {
		return nil, errors.New("Recipient is not valid")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	b, err := json.Marshal(properties)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"issue\",\"contractPayload\":{\"symbol\":\"%s\",\"to\":\"%s\",\"toType\":\"user\",\"feeSymbol\":\"BEE\",\"properties\":%v}}", symbol, to, string(b))
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) TransferNFT(fromName, scid, to, fee string, nfts []api.NFTTransferRequest) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(nfts) <= 0 {
		return nil, errors.New("There is no nft to transfer")
	}
	if len(to) <= 0 {
		return nil, errors.New("Recipient is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	b, err := json.Marshal(nfts)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	data := string(b)
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"transfer\",\"contractPayload\":{\"to\":\"%s\",\"nfts\":%v}}", to, data)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) AddAuthorizedIssuingAccounts(fromName, scid, symbol, fee string, issuingAccounts []string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(issuingAccounts) <= 0 {
		return nil, errors.New("There is no account to add")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	accounts := ""
	for i, item := range issuingAccounts {
		if i == 0 {
			accounts += "[\"" + item + "\""
			if i == len(issuingAccounts)-1 {
				accounts += "]"
			}
		} else if i == len(issuingAccounts)-1 {
			accounts += "," + "\"" + item + "\"]"
		} else {
			accounts += "," + "\"" + item + "\""
		}
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"addAuthorizedIssuingAccounts\",\"contractPayload\":{\"symbol\":\"%s\",\"accounts\":%s}}", symbol, accounts)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) RemoveAuthorizedIssuingAccounts(fromName, scid, symbol, fee string, issuingAccounts []string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(issuingAccounts) <= 0 {
		return nil, errors.New("There is no account to remove")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	accounts := ""
	for i, item := range issuingAccounts {
		if i == 0 {
			accounts += "[\"" + item + "\""
			if i == len(issuingAccounts)-1 {
				accounts += "]"
			}
		} else if i == len(issuingAccounts)-1 {
			accounts += "," + "\"" + item + "\"]"
		} else {
			accounts += "," + "\"" + item + "\""
		}
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"removeAuthorizedIssuingAccounts\",\"contractPayload\":{\"symbol\":\"%s\",\"accounts\":%s}}", symbol, accounts)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) UpdatePropertyDefinition(fromName, scid, symbol, propertyName, newPropertyName, newPropertyType, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(propertyName) <= 0 || len(newPropertyName) <= 0 || len(newPropertyType) <= 0 {
		return nil, errors.New("Property is not valid")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Symbol is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)

	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"updatePropertyDefinition\",\"contractPayload\":{\"symbol\":\"%s\",\"name\":\"%s\",\"type\":\"%s\",\"newName\":\"%s\"}}", symbol, propertyName, newPropertyType, newPropertyName)
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) SetProperties(fromName, scid, symbol, fee string, nfts []api.NFTProperty) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(nfts) <= 0 {
		return nil, errors.New("There is no property to set")
	}
	if len(symbol) <= 0 {
		return nil, errors.New("Recipient is not valid")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	b, err := json.Marshal(nfts)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"setProperties\",\"contractPayload\":{\"symbol\":\"%s\",\"nfts\":%v}}", symbol, string(b))
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) BurnNFT(fromName, scid, fee string, nfts []api.NFTTransferRequest) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(nfts) <= 0 {
		return nil, errors.New("There is no nft to burn")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	b, err := json.Marshal(nfts)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"burn\",\"contractPayload\":{\"nfts\":%v}}", string(b))
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

func (client *Client) MultipleIssueNFT(fromName, scid, fee string, instances []api.Instance) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(instances) <= 0 {
		return nil, errors.New("There is no nft to issue")
	}
	if len(scid) <= 0 {
		scid = "s01"
	}
	var owners []string
	owners = append(owners, fromName)
	b, err := json.Marshal(instances)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	var scoperation string
	scoperation = fmt.Sprintf("{\"contractName\":\"nft\",\"contractAction\":\"issueMultiple\",\"contractPayload\":{\"instances\":%v}}", string(b))
	var trx []types.Operation
	tx := &types.SmartContractOperation{
		RequiredOwners: owners,
		Scid:           scid,
		ScOperation:    scoperation,
		Fee:            fee,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmartContract", Bresp: resp}, err
}

//Transfer of funds to any user.
func (client *Client) Transfer(fromName, toName, memo, amount, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	validate = ValidateAmount(amount)
	if validate == false {
		return nil, errors.New("Amount is not valid")
	}
	var trx []types.Operation
	tx := &types.TransferOperation{
		From:   fromName,
		To:     toName,
		Amount: amount,
		Fee:    fee,
		Memo:   memo,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "Transfer", Bresp: resp}, err
}

func (client *Client) TransferEx(fromName, toName, memo, amount, fee string, extension string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	validate = ValidateAmount(amount)
	if validate == false {
		return nil, errors.New("Amount is not valid")
	}
	var trx []types.Operation
	tx := &types.TransferOperation{
		From:   fromName,
		To:     toName,
		Amount: amount,
		Fee:    fee,
		Memo:   memo,
	}
	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, extension)
	return &OperResp{NameOper: "Transfer", Bresp: resp}, err
}

func (client *Client) MultiOp(trx []types.Operation, extension string) (*OperResp, error) {
	resp, err := client.SendTrx(trx, extension)
	return &OperResp{NameOper: "Multi", Bresp: resp}, err
}

func (client *Client) CreateToken(creator, controlAcc, tokenName string, decimals uint8, maxSupply uint64) (*OperResp, error) {
	var trx []types.Operation
	tx := &types.SmtCreateOperation{
		ControlAccount: controlAcc,
		Symbol:         &types.AssetSymbol{Decimals: decimals, AssetName: tokenName},
		Creator:        creator,
		SmtCreationFee: config.SMT_CREATION_FEE,
		Precision:      decimals,
		Extensions:     [][]interface{}{},
		MaxSupply:      maxSupply,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SmtCreate", Bresp: resp}, err
}

//AccountSupernodeVote of voting for the delegate.
func (client *Client) AccountSupernodeVote(username, supernodeName, fee string, votes int64) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if votes <= 0 {
		return nil, errors.New("Vote is not valid")
	}
	var trx []types.Operation
	tx := &types.AccountSupernodeVoteOperation{
		Account:   username,
		Supernode: supernodeName,
		Approve:   true,
		Votes:     votes,
		Fee:       fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountSupernodeVote", Bresp: resp}, err
}

//Unvote
func (client *Client) AccountSupernodeUnvote(username, supernodeName, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)

	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var trx []types.Operation
	tx := &types.AccountSupernodeVoteOperation{
		Account:   username,
		Supernode: supernodeName,
		Approve:   false,
		Votes:     0,
		Fee:       fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountSupernodeVote", Bresp: resp}, err
}

//TransferToVesting transfer to POWER
func (client *Client) TransferToVesting(from, to, amount, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	validate = ValidateAmount(amount)
	if validate == false {
		return nil, errors.New("Amount is not valid")
	}
	var trx []types.Operation
	tx := &types.TransferToVestingOperation{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "TransferToVesting", Bresp: resp}, err
}

//WithdrawVesting down POWER
func (client *Client) WithdrawVesting(account, vshares, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)

	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	validate = ValidateAmount(vshares)
	if validate == false {
		return nil, errors.New("Amount is not valid")
	}
	var trx []types.Operation
	tx := &types.WithdrawVestingOperation{
		Account:       account,
		VestingShares: vshares,
		Fee:           fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "WithdrawVesting", Bresp: resp}, err
}

//SupernodeUpdate updating delegate data
func (client *Client) SupernodeUpdate(owner, blocksigningkey, fee string) (*OperResp, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)

	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var trx []types.Operation
	tx := &types.SupernodeUpdateOperation{
		Owner:           owner,
		BlockSigningKey: blocksigningkey,
		Fee:             fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "SupernodeUpdate", Bresp: resp}, err
}

//AccountCreate creating a user in systems
func (client *Client) GenKeys(newAccountName string) (*WalletData, error) {
	role := "owner"
	password, err := RandStringBytes(128)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}
	priv := CreatePrivateKey(newAccountName, role, password)
	pub := CreatePublicKey(config.ADDRESS_PREFIX, priv)

	return &WalletData{Name: newAccountName, PrivateKey: priv, PublicKey: pub}, nil
}

func (client *Client) AccountCreate(creator, newAccountName, publicKey, fee string) (*OperResp, error) {
	err := ValidateNameAccount(newAccountName)
	if err != nil {
		return nil, err
	}
	validate := ValidateFee(fee, config.MIN_ACCOUNT_CREATION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	var trx []types.Operation
	empty := map[string]int64{}

	owner := types.Authority{
		WeightThreshold: 1,
		AccountAuths:    empty,
		KeyAuths:        map[string]int64{publicKey: 1},
	}

	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountCreate", Bresp: resp}, err
}

//AccountUpdate update public key for account
func (client *Client) AccountUpdate(account, publicKey, fee string) (*OperResp, error) {
	err := ValidateNameAccount(account)
	if err != nil {
		return nil, err
	}
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}

	var trx []types.Operation
	empty := map[string]int64{}

	owner := types.Authority{
		WeightThreshold: 1,
		AccountAuths:    empty,
		KeyAuths:        map[string]int64{publicKey: 1},
	}
	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountUpdateOperation{
		Account:      account,
		Owner:        &owner,
		JSONMetadata: jsonMeta,
		Fee:          fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountUpdate", Bresp: resp}, err
}

func (client *Client) AccountCreateWS(creator, newAccountName, password, fee string) (*OperResp, error) {
	err := ValidateNameAccount(newAccountName)
	if err != nil {
		return nil, err
	}
	validate := ValidateFee(fee, config.MIN_ACCOUNT_CREATION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	type Keys struct {
		Private string
		Public  string
	}

	var trx []types.Operation
	var listKeys = make(map[string]Keys)
	empty := map[string]int64{}
	roles := [1]string{"owner"}

	for _, val := range roles {
		priv := CreatePrivateKey(newAccountName, val, password)
		pub := CreatePublicKey(config.ADDRESS_PREFIX, priv)
		listKeys[val] = Keys{Private: priv, Public: pub}
	}

	owner := types.Authority{
		WeightThreshold: 1,
		AccountAuths:    empty,
		KeyAuths:        map[string]int64{listKeys["owner"].Public: 1},
	}

	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountCreateWS", Bresp: resp}, err
}

//CreateMultiSigAccount creating an account shared among many users in systems
func (client *Client) CreateMultiSigAccount(creator, newAccountName, fee string, accountOwners []string, keyOwners []string,
	threshold uint32) (*OperResp, error) {
	err := ValidateNameAccount(newAccountName)
	if err != nil {
		return nil, err
	}
	validate := ValidateFee(fee, config.MIN_ACCOUNT_CREATION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(keyOwners)+len(accountOwners) == 0 {
		return nil, errors.New("accountOwners + keyOwners is not empty")
	}
	if threshold == 0 || threshold > uint32(len(keyOwners)+len(accountOwners)) {
		return nil, errors.New("threshold is not valid")
	}
	//Sort owners
	if len(accountOwners) > 1 {
		sort.Strings(accountOwners)
	}
	if len(keyOwners) > 1 {
		sort.Strings(keyOwners)
	}

	var trx []types.Operation
	var mapAccount = make(map[string]int64)
	for _, k := range accountOwners {
		mapAccount[k] = 1
	}
	var mapKeys = make(map[string]int64)
	for _, k := range keyOwners {
		mapKeys[k] = 1
	}

	owner := types.Authority{
		WeightThreshold: threshold,
		AccountAuths:    mapAccount,
		KeyAuths:        mapKeys,
	}
	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountCreateOperation{
		Fee:            fee,
		Creator:        creator,
		NewAccountName: newAccountName,
		Owner:          &owner,
		JSONMetadata:   jsonMeta,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountCreate", Bresp: resp}, err
}

//AccountUpdate update owner keys for account
//TODO: every key has different weight on account
func (client *Client) UpdateMultiSigAccount(account, fee string, accountOwners []string, keyOwners []string, threshold uint32) (*OperResp, error) {
	err := ValidateNameAccount(account)
	if err != nil {
		return nil, err
	}
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	if len(keyOwners)+len(accountOwners) == 0 {
		return nil, errors.New("accountOwners + keyOwners is not empty")
	}
	if threshold == 0 || threshold > uint32(len(keyOwners)+len(accountOwners)) {
		return nil, errors.New("threshold is not valid")
	}
	//Sort owners
	if len(accountOwners) > 1 {
		sort.Strings(accountOwners)
	}
	if len(keyOwners) > 1 {
		sort.Strings(keyOwners)
	}

	var trx []types.Operation
	var mapAccount = make(map[string]int64)
	for _, k := range accountOwners {
		mapAccount[k] = 1
	}
	var mapKeys = make(map[string]int64)
	for _, k := range keyOwners {
		mapKeys[k] = 1
	}

	owner := types.Authority{
		WeightThreshold: threshold,
		AccountAuths:    mapAccount,
		KeyAuths:        mapKeys,
	}
	jsonMeta := &types.AccountMetadata{}
	tx := &types.AccountUpdateOperation{
		Account:      account,
		Owner:        &owner,
		JSONMetadata: jsonMeta,
		Fee:          fee,
	}

	trx = append(trx, tx)
	resp, err := client.SendTrx(trx, "")
	return &OperResp{NameOper: "AccountUpdate", Bresp: resp}, err
}

func (client *Client) CreateTrxTransfer(fromName, toName, memo, amount, fee string, extension string) (*transactions.SignedTransaction, error) {
	validate := ValidateFee(fee, config.MIN_TRANSACTION_FEE)
	if validate == false {
		return nil, errors.New("Fee is not valid")
	}
	validate = ValidateAmount(amount)
	if validate == false {
		return nil, errors.New("Amount is not valid")
	}
	var trxOps []types.Operation
	tOp := &types.TransferOperation{
		From:   fromName,
		To:     toName,
		Amount: amount,
		Fee:    fee,
		Memo:   memo,
	}
	trxOps = append(trxOps, tOp)

	// CreateTrx
	tx, err := client.CreateTrx(trxOps, extension)

	return tx, err
}

func (client *Client) CreateTrx(trxOps []types.Operation, extension string) (*transactions.SignedTransaction, error) {
	// Getting the necessary parameters
	refBlockNum, err := client.GetHeadBlockNum()
	if err != nil {
		return nil, err
	}
	block, err := client.API.GetBlock(refBlockNum)
	if err != nil {
		return nil, err
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return nil, err
	}

	ex := []interface{}{}
	if len(extension) > 0 {
		ex = make([]interface{}, 1)
		as := types.ExtensionJsonType{extension}
		tas := types.ExtensionType{uint8(types.ExtJsonType.Code()), as}
		ex[0] = &tas
	}

	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(refBlockNum),
		RefBlockPrefix: refBlockPrefix,
		Extensions:     ex,
	})

	// Adding Operations to a Transaction
	for _, val := range trxOps {
		tx.PushOperation(val)
	}

	expTime := time.Now().Add(config.TRANSACTION_EXPIRATION_IN_MIN * time.Minute).UTC()
	tm := types.Time{
		Time: &expTime,
	}
	tx.Expiration = &tm

	createdTime := time.Now().UTC()
	tx.CreatedTime = types.UInt64(createdTime.Unix())

	tx.Transaction.Signatures = []string{}

	return tx, err
}

func (client *Client) SignTrx(tx *transactions.SignedTransaction) (*transactions.SignedTransaction, error) {
	// Obtain the key required for signing
	privKeys, err := client.GetSigningKeysOwner()
	if err != nil {
		return nil, err
	}

	// Sign the transaction
	txId, err := tx.Sign(privKeys, client.chainID)
	if err != nil || txId == "" {
		return nil, err
	}
	return tx, nil
}

func (client *Client) SignTrxMulti(tx *transactions.SignedTransaction) ([]string, error) {
	var sigsHex []string
	// Obtain the key required for signing
	privKeys, err := client.GetSigningKeysOwner()
	if err != nil {
		return sigsHex, err
	}

	// Sign the transaction
	sigsHex, err = tx.SignMulti(privKeys, client.chainID)
	if err != nil {
		return sigsHex, err
	}
	return sigsHex, nil
}

func (client *Client) SendTrxMultiSig(tx *transactions.SignedTransaction) (*BResp, error) {
	var bresp BResp
	var err error
	// Sending a transaction
	if client.AsyncProtocol {
		var resp *api.AsyncBroadcastResponse
		resp, err = client.API.BroadcastTransaction(tx.Transaction)
		if resp != nil {
			bresp.ID = resp.ID
		}
	} else {
		var resp *api.BroadcastResponse
		resp, err = client.API.BroadcastTransactionSynchronous(tx.Transaction)
		if resp != nil {
			bresp.ID = resp.ID
		}
	}

	bresp.JSONTrx, _ = JSONTrxString(tx)

	if err != nil {
		return &bresp, err
	}

	return &bresp, nil
}

func ValidateNameAccount(name string) error {
	if name == "" {
		return errors.New("Name account is not empty")
	}
	if len(name) < 3 || len(name) > 16 {
		return errors.New("Name length is from 3 to 16 characters")
	}
	for _, c := range name {
		if !strings.Contains(config.NAME_LETTER, string(c)) {
			return errors.New("Name contains character invalid")
		}
	}
	return nil
}

func ValidateFee(fee string, minFee float64) bool {
	//Validate format of fee
	amount := strings.Split(fee, " ")[0]
	symbol := strings.Split(fee, " ")[1]
	if symbol != config.WD_SYMBOL {
		return false
	}
	//Parse fee to float
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false
	}
	if amt < minFee {
		return false
	}
	return true
}

func ValidateAmount(amount string) bool {
	amtStr := strings.Split(amount, " ")[0]
	amt, err := strconv.ParseFloat(amtStr, 64)
	if err != nil {
		return false
	}
	if amt <= 0 {
		return false
	}
	return true
}

func ValidateNftSymbol(symbol string) error {
	if symbol == "" {
		return errors.New("Symbol is not empty")
	}
	if len(symbol) > 100 {
		return errors.New("Symbol length is maximum 100 characters")
	}
	for _, c := range symbol {
		if !strings.Contains(config.SYMBOL_LETTER, string(c)) {
			return errors.New("Symbol contains character invalid")
		}
	}
	return nil
}
