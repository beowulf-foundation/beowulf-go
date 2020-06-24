package client

import (
	"beowulf-go/api"
	"beowulf-go/config"
	"beowulf-go/transactions"
	"beowulf-go/types"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
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
	password := RandStringBytes(128)
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
