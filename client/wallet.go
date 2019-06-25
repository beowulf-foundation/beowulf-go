package client

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"os"
)

var (
	Keys_ = make(map[string]string)
	Checksum_ = make([]byte, sha512.Size)
	Wallet_ Wallet
	Locked = true
	WalletName_ = "wallet.json"
)

const letterBytes = "0123456789+/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

type PlainKeys struct{
	Checksum [sha512.Size]byte	`json:"checksum"`
	Keys map[string]string 		`json:"keys"`
}

type Wallet struct{
	CipherKeys string `json:"cipher_keys"`
	CipherType string `json:"cipher_type"` // "aes-256-cbc"
	Salt 	   string `json:"salt"`
	Name	   string `json:"name"`
}

func isNew() bool {
	return len(Wallet_.CipherKeys) == 0
}

func isLocked() bool {
	return allZero(Checksum_)
}

func lock() error{
	if(isLocked()){
		return errors.New("The wallet must be unlocked before the password can be set")
	}
	err := encryptKeys();
	if err != nil{
		return err
	}
	for k := range Keys_ {
		Keys_[k] = ""
	}
	Keys_ = make(map[string]string)
	Checksum_ = Checksum_[:0]
	Locked = true
	return nil
}

func (client *Client) Unlock(password string) error {
	if(len(password) == 0){
		return errors.New("Password must be not empty")
	}
	new_password := password + Wallet_.Salt
	pw := sha512.Sum512([]byte(new_password))
	decrypted, err := Decrypt(pw[:], Wallet_.CipherKeys)
	if err != nil{
		return err
	}
	var pk PlainKeys
	err = json.Unmarshal([]byte(decrypted), &pk)
	if err != nil{
		return err
	}
	if(pw != pk.Checksum){
		return errors.New("Don't match checksum")
	}
	Keys_ = pk.Keys
	Checksum_ = pk.Checksum[:]
	Locked = false
	//Set keys
	for k := range pk.Keys {
		client.SetKeys(&Keys{OKey: []string{pk.Keys[k]}})
	}
	return nil
}

func (client *Client) SetPassword(password string) error{
	if (!isNew()) {
		if(isLocked()){
			return errors.New("The wallet must be unlocked before the password can be set")
		}
	}
	salt := randStringBytes(16);
	new_password := password + salt;
	c := sha512.Sum512([]byte(new_password))
	Checksum_ = c[:]
	Wallet_.Salt = salt
	Wallet_.CipherType = "aes-256-cbc"
	return lock()
}

func (client *Client) LoadWallet(wallet_filename string) bool{
	if (wallet_filename == "") {
		wallet_filename = WalletName_
	}
	if _, err := os.Stat(wallet_filename); os.IsNotExist(err) {
		// wallet_filename does not exist
		return false
	}

	dat, err := ioutil.ReadFile(wallet_filename)
	if(dat == nil || err != nil){
		return false
	}
	err = json.Unmarshal(dat, &Wallet_)
	if err != nil{
		return false
	}
	fmt.Println(Wallet_)
	return true;
}

func saveWallet(wallet_filename string) error{
//
// Serialize in memory, then save to disk
//
// This approach lessens the risk of a partially written wallet
// if exceptions are thrown in serialization
//
	err := encryptKeys();
	if err != nil{
		return err
	}

	if (wallet_filename == "") {
		wallet_filename = WalletName_
	}

	data, err := json.Marshal(Wallet_)
	if(err == nil){
		f, _ := os.OpenFile(wallet_filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
		defer f.Close()
		f.Write(data)
	}
	return err
}

func import_key(wif_key, prefix string) bool{
	pubKey := GetPublicKey(prefix, wif_key)
	Keys_[pubKey] = wif_key
	return true
}

func (client *Client) ImportKey(wif_key, name string) bool{
	if(isLocked()){
		return false
	}
	Wallet_.Name = name;

	config, err := client.API.GetConfig()
	if err != nil {
		return false
	}
	prefix := config.AddressPrefix

	if (import_key(wif_key, prefix)) {
		err := saveWallet(name+".json")
		if err != nil{
			return false
		}
		client.SetKeys(&Keys{OKey: []string{wif_key}})
		return true
	}
	return false
}

func encryptKeys() error {
	var data PlainKeys
	data.Keys = Keys_
	copy(data.Checksum[:], Checksum_[:])
	//data.Checksum = string(Checksum_)
	plainData, err := json.Marshal(data)
	if(err != nil){
		return err
	}
	plainTxt := string(plainData)
	Wallet_.CipherKeys, err = Encrypt(data.Checksum[:], plainTxt)
	if err != nil{
		return err
	}
	return nil
}