# beowulf-go

beowulf-go is the official Beowulf library for Go.  

## Usage

```go
import "github.com/beowulf-foundation/beowulf-go"
```

## Example

```go
//1. Init
key := "5Jxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Replace your private key
url := " http://localhost:8376" // Replace this url with your node url
cls, _ := client.NewClient(url)
cls.SetKeys(&client.Keys{OKey: []string{key}})
defer cls.Close()

//2. Get config
log.Println("---> GetConfig()")
config, err := cls.API.GetConfig()
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
json_trx, _ := json.Marshal(trx)
fmt.Println(string(json_trx))

//6. Create account
resp_ac, _ := cls.AccountCreate("creator", "new_account_name","password","1.00000 W")
json_rac, _ := json.Marshal(resp_ac)
fmt.Println(string(json_rac))

//7. Transfer native coin
//7.1. Transfer BWF from alice to bob
resp_bwf, _ := cls.Transfer("alice", "bob", "", "100.00000 BWF", "0.01000 W")
json_rbwf, _ := json.Marshal(resp_bwf)
fmt.Println(string(json_rbwf))

//7.2. Transfer W from alice to bob
resp_w, _ := cls.Transfer("alice", "bob", "", "10.00000 W", "0.01000 W")
json_rw, _ := json.Marshal(resp_w)
fmt.Println(string(json_rw))

//8. Transfer token
//Transfer token KNOW from alice to bob
resp_tk, _ := cls.Transfer("alice", "bob", "", "1000.00000 KNOW", "0.01000 W")
json_rtk, _ := json.Marshal(resp_tk)
fmt.Println(string(json_rtk))
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
