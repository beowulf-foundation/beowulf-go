# beowulf-go

beowulf-go is the official Beowulf library for Go.  

## Install
```go
go get -u github.com/beowulf-foundation/beowulf-go
```

## Usage

```go
import "github.com/beowulf-foundation/beowulf-go"
```

## Example

```go
//1. Init
//// MainNet: https://bw.beowulfchain.com/rpc
//// TestNet: https://testnet-bw.beowulfchain.com/rpc
url := " http://localhost:8376" // Replace this url with your node url
cls, _ := client.NewClient(url)
defer cls.Close()
//// SetKeys
key := "5Jxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Replace your private key
cls.SetKeys(&client.Keys{OKey: []string{key}})


//2. Get config
fmt.Println("========== GetConfig ==========")
config, err := cls.API.GetConfig()
if err != nil {
    fmt.Println(err)
}
json_cfg, _ := json.Marshal(config)
fmt.Println(string(json_cfg))

// Use the last irreversible block number as the initial last block number.
props, err := cls.API.GetDynamicGlobalProperties()
json_props, _ := json.Marshal(props)
fmt.Println(string(json_props))

//3. Get account
account, err := cls.GetAccount("name-account")
json_acc, _ := json.Marshal(account)
fmt.Println(string(json_acc))

//4. Get block
lastBlock := props.LastIrreversibleBlockNum
block, err := cls.GetBlock(lastBlock)
json_bk, _ := json.Marshal(block)
fmt.Println(string(json_bk))

//5. Get transaction
trx, err := cls.API.GetTransaction("673fbd4609d1156bcf6d9e6c36388926f7116acc")
if err != nil {
    fmt.Println(err)
}
json_trx, _ := json.Marshal(trx)
fmt.Println(string(json_trx))
oplist := *trx.Operations
for _, op := range oplist {
    d := op.Data()
    switch d.(type){
    case *types.TransferOperation:
        byteData, _ := json.Marshal(d)
        oop := types.TransferOperation{}
        json.Unmarshal(byteData, &oop)
        fmt.Println(oop)
        fmt.Println("From:", oop.From)
        fmt.Println("To:", oop.To)
        fmt.Println("Amount:", oop.Amount)
        fmt.Println("Fee:", oop.Fee)
        fmt.Println("Memo:", oop.Memo)
    }
}
exlist := trx.Extensions
if len(exlist) > 0 {
    tmp := exlist[0]
    byteex, _ := json.Marshal(tmp)
    var met map[string]interface{}
    json.Unmarshal(byteex, &met)

    et := types.ExtensionType{}
    stype := fmt.Sprintf("%v", met["type"])
    et.Type = uint8(types.GetExtCodes(stype))

    value := met["value"].(map[string]interface{})
    ejt := types.ExtensionJsonType{}
    ejt.Data = fmt.Sprintf("%v", value["data"])
    et.Value = ejt

    fmt.Println(ejt)
    fmt.Println(et)
}

//6. Transfer native coin
//6.1. Transfer BWF from alice to bob
resp_bwf, err := cls.Transfer("alice", "bob", "", "100.00000 BWF", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rbwf, _ := json.Marshal(resp_bwf)
fmt.Println(string(json_rbwf))

//6.2. Transfer W from alice to bob
resp_w, err := cls.Transfer("alice", "bob", "", "10.00000 W", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rw, _ := json.Marshal(resp_w)
fmt.Println(string(json_rw))

//7. Transfer token
//Transfer token KNOW from alice to bob
resp_tk, err := cls.Transfer("alice", "bob", "", "1000.00000 KNOW", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rtk, _ := json.Marshal(resp_tk)
fmt.Println(string(json_rtk))


//8. Create account

//// 8.1. GenKeys
walletData, _ := cls.GenKeys("new-account-name")
json_wd, _ := json.Marshal(walletData)
fmt.Println(string(json_wd))

//// 8.2. AccountCreate
resp_ac, err := cls.AccountCreate("creator", walletData.Name, walletData.PublicKey,"1.00000 W")
if err != nil {
    fmt.Println(err)
}
json_rac, _ := json.Marshal(resp_ac)
fmt.Println(string(json_rac))

//// 8.3. Write file wallet.
password := "your_password"
err := client.SaveWalletFile("/path/to/folder/save/wallet", "", password, walletData)
if err != nil {
    fmt.Println(err)
}

//// 8.4. Load file wallet.
rs := cls.SetKeysFromFileWallet("/path/to/folder/save/wallet/new-account-name-wallet.json", password)
if rs != nil {
    fmt.Println(rs)
}
// print keys
fmt.Println(cls.GetPrivateKey())
fmt.Println(cls.GetPublicKey())
```

## Package Organisation

You need to create a `Client` object to be able to do anything.
Then you just need to call `NewClient()`.

Once you create a `Client` object, you can start calling the methods exported
via `beowulfd`'s RPC endpoint by invoking associated methods on the client object.

When looking for a method to call, all you need is to turn the method name into
CamelCase, e.g. `get_config` becomes `Client.API.GetConfig`.

## License

MIT, see the `LICENSE` file.
