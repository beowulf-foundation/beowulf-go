package main

import (
	"beowulf-go/client"
	"beowulf-go/transactions"
	"beowulf-go/util"
	"encoding/json"
	"fmt"
	"github.com/shettyh/threadpool"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
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
const url string = "https://testnet-bw.beowulfchain.com/rpc" // Replace this url with your node url
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
	CheckNFT()
	//CreateNFT()
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

	nfts, err := cli.GetNFTs("")
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("NFTs", nfts)
	}
}
