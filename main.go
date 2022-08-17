package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/beowulf-foundation/beowulf-go/api"
	"github.com/beowulf-foundation/beowulf-go/client"
	"github.com/beowulf-foundation/beowulf-go/transactions"
	"github.com/beowulf-foundation/beowulf-go/util"
	"github.com/shettyh/threadpool"
)

// 1. https://gobyexample.com/worker-pools
// 2. https://github.com/shettyh/threadpool
// 3. https://medium.com/@j.d.livni/write-a-go-worker-pool-in-15-minutes-c9b42f640923
// 4. https://brandur.org/go-worker-pool
// 5. https://github.com/ivpusic/grpool

const N_ACC int32 = 100
const NI int = int(N_ACC)

var MAP_CLIENT = make(map[string]*client.Client)
var MAP_NAME = make(map[int32]string)

const IsAsync bool = true
const work_dir = "/home/nghiatc/go-projects/src/beowulf-go"
const file_path = work_dir + string(os.PathSeparator) + "bwc100.json"
const url string = "http://35.220.140.28:8376"

// const url string = "https://testnet-bw.beowulfchain.com/rpc" // Replace this url with your node url
const creator = "nghia"

//const key = "5JEUvsDmUxLhPZLsQdofQqMnFxBob6bpXcmJ5i54stf34PqQfyb" // nghia
const key = "5JHTf7dkpVxQNcb5NWc7URTrHDgAFEyxn2BEnMjuJ6fJrCAniCQ" // beowulf
var BLOCK_NUM uint32 = 0
var REF_BLOCK_PREFIX uint32 = 0

const BWF string = "BWF"
const W string = "W"
const NumWorker int = 100
const NumTx int = 100

//var countTx uint64 = 0
//var poolR = threadpool.NewThreadPool(NumWorker, int64(NumWorker * NumTx + 1000))
var poolC = threadpool.NewThreadPool(NumWorker, int64(NumWorker*NumTx+1000))
var mapFTask = make(map[int32]*threadpool.Future)

func main() {
	//// 0. Init
	Init()
	defer Close()
	//GetRefBlock()

	/*
		NFT testing
	*/
	//CheckNFT()
	CreateNFT()
	CommitSidechainBlock()
	//UpdateUrlAndImage()
	//UpdateName()
	//UpdateOrgName()
	AddProperty("name", "string")
	AddProperty("image", "string")
	AddProperty("externalurl", "string")
	AddProperty("description", "string")
	IssueNFT()
	//IssueWithProperty()
	IssueWithProperty2()
	//AddAuthorizedIssuingAccounts()
	//RemoveAuthorizedIssuingAccounts()
	//UpdatePropertyDefinition()
	//SetProperties()
	//MultipleIssueNFT()
	TransferNFT()
	BurnNFT()

	//GetNFTBalance()
	//GetNFTInstances()
	//GetNFTBalanceOfAccount()
	GetNFTTransaction()
	GetLatestNFTBlock()
	GetNFTBlock()

	//// 1. ZeroCoin && IssueCoin
	//CheckTotalCoin()
	//ZeroCoin()

	//IssueCoin()
	//CheckCoin()
	//CheckTotalCoin()

	//// 2. RunPoolNTask
	//CheckTotalCoin()
	//RunPoolNTask()
	//CheckTotalCoin()

	//CheckTotalCoin()
	//for i := 0; i < NumTx; i++ {
	//	from,to := RandFromTo()
	//	TransferBWF(from, to)
	//	atomic.AddUint64(&countTx, 1)
	//}
	//CheckTotalCoin()

	//time.Sleep(2 * time.Second)
	//countFinal := atomic.LoadUint64(&countTx)
	//fmt.Println("========>> countTx:", countFinal)

	//CheckTotalCoin()
	//CheckCoin()

	//CheckTotalCoin()

	time.Sleep(1 * time.Second)
	fmt.Println("====================================================================")
	fmt.Println("==============================End Main==============================")
	fmt.Println("====================================================================")
}

func CommitSidechainBlock() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})
	c := fmt.Sprintf("0x58e983428796f23f235ea13bf02039030464d947262b39fef842c8f97eee8ee7125")
	operResp, err := cli.CommitScb("edge", "beowulf", c, "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(operResp)
		fmt.Println("Tx id: ", operResp.Bresp.ID)
		time.Sleep(5 * time.Second)
		tx, er := cli.GetTransaction(operResp.Bresp.ID)
		if er != nil {
			fmt.Println(er)
		} else {
			fmt.Println("Result tx: \n", tx)
		}
	}
}

func CreateNFT() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	accounts := []string{}
	res, err := cli.CreateNFT("beowulf", "s01", "GONFT", "GONFT", "10000000", "0.01000 W", accounts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

//func UpdateUrl() {
//	cli, _ := client.NewClient(url, true)
//	defer cli.Close()
//	cli.SetKeys(&client.Keys{OKey: []string{key}})
//
//	res, err := cli.UpdateUrl("beowulf", "s01", "GONFT", "https://google.com", "0.01000 W")
//	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(res)
//	}
//}
//
//func UpdateImage() {
//	cli, _ := client.NewClient(url, true)
//	defer cli.Close()
//	cli.SetKeys(&client.Keys{OKey: []string{key}})
//
//	res, err := cli.UpdateImage("beowulf", "s01", "GONFT", "https://image.com", "0.01000 W")
//	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(res)
//	}
//}

func UpdateUrlAndImage() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	res, err := cli.UpdateMetadata("beowulf", "s01", "GONFT", "https://url.com", "https://image.update.vn", "0.01000 W")
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func UpdateName() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	res, err := cli.UpdateName("beowulf", "s01", "GONFT", "HELLONFT", "0.01000 W")
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func UpdateOrgName() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	res, err := cli.UpdateOrgName("beowulf", "s01", "GONFT", "ORG", "0.01000 W")
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func AddProperty(propertyName, propertyType string) {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	accounts := []string{}
	accounts = append(accounts, "beowulf")
	res, err := cli.AddProperty("beowulf", "s01", "GONFT", propertyName, propertyType, "0.01000 W", accounts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func IssueNFT() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	res, err := cli.IssueNFT("beowulf", "s01", "GONFT", "beowulf", "0.01000 W")
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

type PropertyRequest struct {
	RoomId string `json:"roomId"`
}

type PropertyRequest2 struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Externalurl string `json:"externalurl"`
	Description string `json:"description"`
}

func IssueWithProperty() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	//var properties PropertyRequest
	//properties.RoomName = "A101"
	properties := &PropertyRequest{RoomId: "001"}
	res, err := cli.IssueWithProperties("beowulf", "s01", "GONFT", "beowulf", "0.01000 W", properties)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func IssueWithProperty2() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	//var properties PropertyRequest
	//properties.RoomName = "A101"
	properties := &PropertyRequest2{Name: "Oriental Nha Trang Hotel: Stay on Jul 02, 2022", Image: "https://crystabaya-public-staging.s3.ap-southeast-1.amazonaws.com/arts/3.jpg", Externalurl: "https://crystabaya.com/property?id=29vXQMeu7tAacNZSvbF3UVSd6Fl", Description: "This Crystabaya NFT token represents one hotel stay at the Oriental Nha Trang Hotel for a stay on Jul 02, 2022.\\n\\nTo redeem this token, please do the following steps:\\n\\n1. Open an account on [Crystabaya.com](https://crystabaya.com)\\n\\n2. Get the wallet address for your account\\n\\n3. Deposit this token into the wallet address above\\n\\nThat's it. You will be able to check-in at the hotel for the stay on the date above.\\n\\nOnce redeemed, you will also receive a hand-drawing collectible artwork included in this NFT token as a souvenir. You can withdraw it back to your own wallet and/or trade it on OpenSea if desirable."}
	res, err := cli.IssueWithProperties("beowulf", "s01", "GONFT", "khoa01", "0.01000 W", properties)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func AddAuthorizedIssuingAccounts() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	accounts := []string{}
	accounts = append(accounts, "beowulf2")
	res, err := cli.AddAuthorizedIssuingAccounts("beowulf", "s01", "GONFT", "0.01000 W", accounts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func RemoveAuthorizedIssuingAccounts() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	accounts := []string{}
	accounts = append(accounts, "beowulf2")
	res, err := cli.RemoveAuthorizedIssuingAccounts("beowulf", "s01", "GONFT", "0.01000 W", accounts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func UpdatePropertyDefinition() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	res, err := cli.UpdatePropertyDefinition("beowulf", "s01", "GONFT", "roomName", "roomCode", "string", "0.01000 W")
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func SetProperties() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	var properties []api.NFTProperty
	var property api.NFTProperty
	property.Id = "2"
	var p api.Property
	p.Name = "roomId"
	p.Data = "206"
	property.Properties = p
	properties = append(properties, property)
	res, err := cli.SetProperties("beowulf", "s01", "GONFT", "0.01000 W", properties)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func MultipleIssueNFT() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	var instances []api.Instance
	var instance1 api.Instance
	instance1.To = "beowulf"
	instance1.Symbol = "GONFT"
	instance1.FeeSymbol = "BEE"
	instance1.ToType = "user"
	var instance2 api.Instance
	instance2.To = "beowulf2"
	instance2.Symbol = "GONFT"
	instance2.FeeSymbol = "BEE"
	instance2.ToType = "user"
	instances = append(instances, instance1, instance2)
	res, err := cli.MultipleIssueNFT("beowulf", "s01", "0.01000 W", instances)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func TransferNFT() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	var nfts []api.NFTTransferRequest
	var nft api.NFTTransferRequest
	var ids []string
	ids = append(ids, "1")
	nft.Symbol = "GONFT"
	nft.Ids = ids
	nfts = append(nfts, nft)
	res, err := cli.TransferNFT("beowulf", "s01", "beowulf2", "0.01000 W", nfts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func BurnNFT() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	var nfts []api.NFTTransferRequest
	var nft api.NFTTransferRequest
	nft.Symbol = "GONFT"
	nft.Ids = append(nft.Ids, "1")
	nfts = append(nfts, nft)
	res, err := cli.BurnNFT("beowulf", "s01", "0.01000 W", nfts)
	//res, err := cli.Transfer("beowulf", "beowulf2", "nothing", "0.10000 W", "0.01000 W")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

/** Pool TaskCallable */
type NTaskCallable struct {
	ID int32
}

func (nc *NTaskCallable) Call() interface{} {
	var count int64 = 0
	fmt.Println("Running NTaskCallable:", nc.ID)
	for i := 0; i < NumTx; i++ {
		from, to := RandFromTo()
		TransferBWF(from, to)
		count++
	}
	fmt.Println("End NTaskCallable:", nc.ID)
	return count
}
func StatPoolNTaskCallable(mapF map[int32]*threadpool.Future) int64 {
	var total int64 = 0
	for k, v := range mapF {
		count := v.Get().(int64)
		fmt.Printf("Count NTaskCallable[%d] = %d\n", k, count)
		total += count
	}
	fmt.Printf("Count StatPoolNTaskCallable = %d\n", total)
	return total
}

/** Pool Task simple */
//type NTask struct {
//	ID int32
//}
//
//func (nt *NTask) Run() {
//	fmt.Println("Running NTask:", nt.ID)
//	for i := 0; i < NumTx; i++ {
//		from,to := RandFromTo()
//		TransferBWF(from, to)
//		atomic.AddUint64(&countTx, 1)
//	}
//	fmt.Println("End NTask:", nt.ID)
//}
//
//func RunPoolNTask() {
//	for id := 0; id < NumWorker; id++ {
//		task := &NTask{ID: int32(id)}
//		poolR.Execute(task)
//	}
//	time.Sleep(1 * time.Second)
//}

func TransferBWF(from int32, to int32) {
	if from != to && 0 <= from && from < N_ACC && 0 <= to && to < N_ACC {
		fromName := MAP_NAME[from]
		toName := MAP_NAME[to]
		fromCli := MAP_CLIENT[fromName]
		memo := "aaaaaa"
		start := time.Now()
		resp_bwf, err := fromCli.Transfer(fromName, toName, memo, "0.01000 BWF", "0.01000 W")
		//resp_bwf, err := fromCli.TransferEx(fromName, toName, memo, "0.01000 BWF", "0.01000 W", BLOCK_NUM, REF_BLOCK_PREFIX)
		period := time.Since(start)
		if err != nil {
			//fmt.Printf("TransferBWF: %v\n", err)
			fmt.Println(err)
		}
		fmt.Println("========== Transfer BWF ==========")
		fmt.Printf("Time Transfer: %d ns\n", period.Nanoseconds())
		// Điều kiện mạng bình thường, thời gian Transfer Async giao động từ: 63 --> 97 ms
		json_rbwf, _ := json.Marshal(resp_bwf)
		fmt.Println(string(json_rbwf))
	} else {
		fmt.Println("*******************From-To invalid*****************")
	}
}

func Init() {
	// 1. Init MapName && MapClient
	var objmap map[string]*client.WalletData
	data, err := ioutil.ReadFile(file_path)
	if data == nil || err != nil {
		fmt.Println(err)
		fmt.Println("File 100 wallet is empty or can not read.")
	}
	if err := json.Unmarshal(data, &objmap); err != nil {
		fmt.Println(err)
	}
	json_omac, _ := json.Marshal(objmap)
	fmt.Println(string(json_omac))
	fmt.Println("===>> Read file 100 wallet complete.")
	for k, v := range objmap {
		cli, _ := client.NewClient(url, true)
		cli.SetKeys(&client.Keys{OKey: []string{v.PrivateKey}})
		cli.AsyncProtocol = IsAsync
		MAP_CLIENT[v.Name] = cli
		nk, _ := strconv.ParseInt(k, 10, 32)
		MAP_NAME[int32(nk)] = v.Name
	}
	json_mcli, _ := json.Marshal(MAP_CLIENT)
	fmt.Println(string(json_mcli))
	json_mname, _ := json.Marshal(MAP_NAME)
	fmt.Println(string(json_mname))

	// 2. Init BLOCK_NUM && REF_BLOCK_PREFIX
	BLOCK_NUM, REF_BLOCK_PREFIX = GetRefBlock()
	fmt.Println("BLOCK_NUM:", BLOCK_NUM)
	fmt.Println("REF_BLOCK_PREFIX:", REF_BLOCK_PREFIX)
}

func Close() {
	for _, c := range MAP_CLIENT {
		c.Close()
	}
	//poolR.Close()
	poolC.Close()
}

func GetRefBlock() (uint32, uint32) {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	blockNum, err := cli.GetHeadBlockNum()
	if err != nil {
		return 0, 0
	}
	block, err := cli.API.GetBlock(blockNum)
	if err != nil {
		return 0, 0
	}
	refBlockId := block.BlockId
	// Creating a Transaction
	refBlockPrefix, err := transactions.RefBlockPrefix(refBlockId)
	if err != nil {
		return 0, 0
	}
	return blockNum, uint32(refBlockPrefix)
}

func RandFromTo() (int32, int32) {
	rand.Seed(time.Now().UnixNano())
	var f = rand.Int31n(N_ACC)
	var t = f
	for {
		t = rand.Int31n(N_ACC)
		if f != t {
			break
		}
	}
	return f, t
}

func IssueCoin() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	for _, name := range MAP_NAME {
		fmt.Println("IssueCoin BWF for: ", name)
		memo := "aaaa"
		_, err := cli.Transfer(creator, name, memo, "1000.00000 BWF", "0.01000 W")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("IssueCoin W for: ", name)
		_, err1 := cli.Transfer(creator, name, memo, "100.00000 W", "0.01000 W")
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}

func ZeroCoin() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()
	cli.SetKeys(&client.Keys{OKey: []string{key}})

	for _, v := range MAP_NAME {
		fmt.Println("ZeroCoin Issue W for: ", v)
		memo := "aaa"
		_, err1 := cli.Transfer(creator, v, memo, "1.00000 W", "0.01000 W")
		if err1 != nil {
			fmt.Println(err1)
		}
	}

	for name, fromCli := range MAP_CLIENT {
		memo := "aaaa"
		f_account, _ := fromCli.GetAccount(name)
		is_bwf := false
		bwf, _ := util.ParseBalance(f_account.Balance)
		if bwf > 0 {
			fmt.Println("ZeroCoin BWF of: ", name)
			_, err := fromCli.Transfer(name, creator, memo, f_account.Balance, "0.01000 W")
			if err != nil {
				fmt.Println(err)
			}
			is_bwf = true
		}
		w, _ := util.ParseBalance(f_account.WdBalance)
		if w > 0 {
			fmt.Println("ZeroCoin W of: ", name)
			if is_bwf {
				w = w - 0.02
			} else {
				w = w - 0.01
			}
			_, err1 := fromCli.Transfer(name, creator, memo, util.FormatBalance(w, W), "0.01000 W")
			if err1 != nil {
				fmt.Println(err1)
			}
		}
	}
}

func CheckCoin() {
	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	name := "bwc"
	count := 0
	for i := 0; i < NI; i++ {
		accName := name + strconv.Itoa(i)
		account, err := cli.GetAccount(accName)
		if err == nil {
			count++
		}
		fmt.Println("========== GetAccount:", accName)
		fmt.Println(account.Balance)
		fmt.Println(account.WdBalance)
	}
	fmt.Println("CheckCoin 100 Account complete:", count)
}

func CheckTotalCoin() (float64, float64) {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	wbfTotal := 0.0
	wTotal := 0.0
	name := "bwc"
	count := 0
	for i := 0; i < NI; i++ {
		accName := name + strconv.Itoa(i)
		account, err := cli.GetAccount(accName)
		if err == nil {
			count++
		}
		bwf, _ := util.ParseBalance(account.Balance)
		wbfTotal += bwf
		w, _ := util.ParseBalance(account.WdBalance)
		wTotal += w
	}
	fmt.Println("CheckTotalCoin 100 Account complete:", count)
	fmt.Println("Total BWF:", util.FormatBalance(wbfTotal, BWF))
	fmt.Println("Total W:", util.FormatBalance(wTotal, W))
	return wbfTotal, wTotal
}

func CheckNFT() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTs("", 100, 0)
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetNFTBalance() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTBalance("khoa02", "CBAYNFT", 100, 0)
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetNFTInstances() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTInstances("CBAYNFT", 100, 0)
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetNFTBalanceOfAccount() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTBalanceOfAccount("khoa02", 100, 0)
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetLatestNFTBlock() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetLatestNFTBlock()
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetNFTBlock() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTBlock(20)
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}

func GetNFTTransaction() {
	time.Sleep(1 * time.Second)

	cli, _ := client.NewClient(url, true)
	defer cli.Close()

	nfts, err := cli.GetNFTTransaction("7129ca688321d386e670b093610bb0db68714eac")
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}
