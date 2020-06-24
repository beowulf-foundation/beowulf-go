# Official Go BEOWULF Library

`beowulf-go` is the official Beowulf library for GoLang.  

## Main Functions Supported
1. CHAIN
- Get block
- Get transaction
- Get account
- Get token
- Get balance
2. TRANSACTION
- Transfer
- Create wallet
- Create token
- Vote for supernode
- Unvote supernode
- Transfer to vesting
- Withdraw

## Installation
```go
go get -u github.com/beowulf-foundation/beowulf-go
```

## Configuration
```go
import (
    "github.com/beowulf-foundation/beowulf-go/api"
    "github.com/beowulf-foundation/beowulf-go/client"
    "github.com/beowulf-foundation/beowulf-go/config"
    "github.com/beowulf-foundation/beowulf-go/encoding"
    "github.com/beowulf-foundation/beowulf-go/transactions"
    "github.com/beowulf-foundation/beowulf-go/transports"
    "github.com/beowulf-foundation/beowulf-go/types"
    "github.com/beowulf-foundation/beowulf-go/util"
)
```
#### Init

```go
// MainNet: https://bw.beowulfchain.com/rpc
// TestNet: https://testnet-bw.beowulfchain.com/rpc
url := "http://localhost:8376/rpc"  #Replace this url with your node url
cls, _ := client.NewClient(url, true)  #Set "true" if using testnet, "false" if using mainnet
defer cls.Close()
// SetKeys
key := "5Jxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Replace your private key
cls.SetKeys(&client.Keys{OKey: []string{key}})
```

## Example Usage
```go
import (
    "github.com/beowulf-foundation/beowulf-go/api"
    "github.com/beowulf-foundation/beowulf-go/client"
    "github.com/beowulf-foundation/beowulf-go/config"
    "github.com/beowulf-foundation/beowulf-go/encoding"
    "github.com/beowulf-foundation/beowulf-go/transactions"
    "github.com/beowulf-foundation/beowulf-go/transports"
    "github.com/beowulf-foundation/beowulf-go/types"
    "github.com/beowulf-foundation/beowulf-go/util"
)
```
##### Get block
```go
lastBlock := props.LastIrreversibleBlockNum
block, err := cls.GetBlock(lastBlock)
json_bk, _ := json.Marshal(block)
fmt.Println(string(json_bk))
```

##### Get transaction
```go
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
```

##### Get account
```go
account := "alice"          #Replace with your account name
result,_ := cls.GetAccount(account)
fmt.Println(result) 
```

##### Get token
```go
token := "KNOW"             #Replace with your token name
result,_ := cls.GetToken(token)
fmt.Println(result) 
```

##### Get balance
```go
account := "alice"          #Replace with your account name
token := "KNOW"             #Replace with your token name
decimal := 5                #Replace with your token decimal
result, _ := cls.API.GetBalance(account, token, decimal)
fmt.Println(result)
```

##### Transfer native coin
###### Transfer BWF
```go
resp_bwf, err := cls.Transfer("alice", "bob", "", "100.00000 BWF", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rbwf, _ := json.Marshal(resp_bwf)
fmt.Println(string(json_rbwf))
```

###### Transfer W
```go
resp_w, err := cls.Transfer("alice", "bob", "", "10.00000 W", "0.01000 W")
if err != nil {
    fmt.Println(err)
}
json_rw, _ := json.Marshal(resp_w)
fmt.Println(string(json_rw))
```

##### Create wallet 

```go
walletData, _ := cls.GenKeys("new-account-name")
json_wd, _ := json.Marshal(walletData)
fmt.Println(string(json_wd))

resp_ac, err := cls.AccountCreate("creator", walletData.Name, walletData.PublicKey,"1.00000 W")
if err != nil {
    fmt.Println(err)
}
json_rac, _ := json.Marshal(resp_ac)
fmt.Println(string(json_rac))

password := "your_password"
err := client.SaveWalletFile("/path/to/folder/save/wallet", "", password, walletData)
if err != nil {
    fmt.Println(err)
}

rs := cls.SetKeysFromFileWallet("/path/to/folder/save/wallet/new-account-name-wallet.json", password)
if rs != nil {
    fmt.Println(rs)
}
// print keys
fmt.Println(cls.GetPrivateKey())
fmt.Println(cls.GetPublicKey())
```

##### Create token

```go
creator := "initminer"      #Account who has authority of creating token
owner := "alice"            #Account who is owner of token
token := "KNOW"             #Token name of token being created
decimal := 5                #Decimal of token being created
maxSupply := 1000000        #Maximum token will be supplied
cls.CreateToken(creator, owner, token, decimal, maxSupply)
```

##### Vote

```go
account := "alice"          #User account go to vote
supernode := "initminer"    #Supernode account is voted 
fee := "0.01000 W"          #Fee to vote 
votes := 100                #Number of vote 
cls.AccountSupernodeVote(account, supernode, fee, votes)
```

##### Unvote 

```go
account := "alice"          #User account go to unvote
supernode := "initminer"    #Supernode account is unvoted  
fee := "0.01000 W"          #Fee to unvote
cls.AccountSupernodeUnvote(account, supernode, fee)
```

##### Transfer to vesting
```go
from := "alice"             #Account who wants to transfer
to := "bob"                 #Receiver
amount := "1000.00000 BWF"  #Amount to transfer
fee := "0.01000 W"          #Fee to transfer
cls.TransferToVesting(from, to, amount, fee)
```

##### Withdraw
```go
account := "alice"          #Account wanting to withdraw vest to BWF
amount := "100.00000 M"     #Number of vest to withdraw
fee := "0.01000 W"          #Fee to withdraw
cls.WithdrawVesting(account, amount, fee)
```
