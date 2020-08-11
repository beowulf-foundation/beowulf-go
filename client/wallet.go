package client

import (
	"beowulf-go/config"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
)

var (
	Keys_       = make(map[string]string)
	Checksum_   = make([]byte, sha512.Size)
	Wallet_     Wallet
	Locked      = true
	WalletName_ = "wallet.json"
)

const letterBytes = "0123456789+abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//func RandStringBytes(n int) string {
//	rand.Seed(time.Now().UnixNano())
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}

func RandStringBytes(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		ret = append(ret, letterBytes[num.Uint64()])
	}

	return string(ret), nil
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

type PlainKeys struct {
	Checksum [sha512.Size]byte `json:"checksum"`
	Keys     map[string]string `json:"keys"`
}

type Wallet struct {
	CipherKeys string `json:"cipher_keys"`
	CipherType string `json:"cipher_type"` // "aes-256-cbc"
	Salt       string `json:"salt"`
	Name       string `json:"name"`
}

type WalletData struct {
	Name       string `json:"name"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func isNew() bool {
	return len(Wallet_.CipherKeys) == 0
}

func isLocked() bool {
	return allZero(Checksum_)
}

func lock() error {
	if isLocked() {
		return errors.New("The wallet must be unlocked before the password can be set")
	}
	err := encryptKeys()
	if err != nil {
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
	if len(password) == 0 {
		return errors.New("Password must be not empty")
	}
	new_password := password + Wallet_.Salt
	pw := sha512.Sum512([]byte(new_password))
	decrypted, err := Decrypt(pw[:], Wallet_.CipherKeys)
	if err != nil {
		return err
	}
	var pk PlainKeys
	err = json.Unmarshal([]byte(decrypted), &pk)
	if err != nil {
		return err
	}
	if pw != pk.Checksum {
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

func (client *Client) SetPassword(password string) error {
	if !isNew() {
		if isLocked() {
			return errors.New("The wallet must be unlocked before the password can be set")
		}
	}
	salt, err := RandStringBytes(16)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}
	new_password := password + salt
	c := sha512.Sum512([]byte(new_password))
	Checksum_ = c[:]
	Wallet_.Salt = salt
	Wallet_.CipherType = "aes-256-cbc"
	return lock()
}

func (client *Client) LoadWallet(wallet_filename string) bool {
	if wallet_filename == "" {
		wallet_filename = WalletName_
	}
	if _, err := os.Stat(wallet_filename); os.IsNotExist(err) {
		// wallet_filename does not exist
		return false
	}

	dat, err := ioutil.ReadFile(wallet_filename)
	if dat == nil || err != nil {
		return false
	}
	err = json.Unmarshal(dat, &Wallet_)
	if err != nil {
		return false
	}
	return true
}

func saveWallet(wallet_filename string) error {
	//
	// Serialize in memory, then save to disk
	//
	// This approach lessens the risk of a partially written wallet
	// if exceptions are thrown in serialization
	//
	err := encryptKeys()
	if err != nil {
		return err
	}

	if wallet_filename == "" {
		wallet_filename = WalletName_
	}

	data, err := json.Marshal(Wallet_)
	if err == nil {
		f, _ := os.OpenFile(wallet_filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
		defer f.Close()
		f.Write(data)
	}
	return err
}

func import_key(wif_key, prefix string) bool {
	pubKey := CreatePublicKey(prefix, wif_key)
	Keys_[pubKey] = wif_key
	return true
}

func (client *Client) ImportKey(wif_key, name string) bool {
	if isLocked() {
		return false
	}
	Wallet_.Name = name

	if import_key(wif_key, config.ADDRESS_PREFIX) {
		err := saveWallet(name + ".json")
		if err != nil {
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
	plainData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	plainTxt := string(plainData)
	Wallet_.CipherKeys, err = Encrypt(data.Checksum[:], plainTxt)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) SetKeysFromFileWallet(pathFileWallet string, password string) error {
	if pathFileWallet == "" {
		return errors.New("Path file wallet is not empty.")
	}
	if password == "" {
		return errors.New("Password is not empty.")
	}

	if _, err := os.Stat(pathFileWallet); os.IsNotExist(err) {
		// pathFileWallet does not exist
		return errors.New("File wallet does not exist.")
	}

	data, err := ioutil.ReadFile(pathFileWallet)
	if data == nil || err != nil {
		return errors.New("File wallet is empty or can not read.")
	}
	var wl *Wallet
	err = json.Unmarshal(data, &wl)
	if wl == nil || err != nil {
		return errors.New("Can not decode json wallet data.")
	}

	new_password := password + wl.Salt
	pw := sha512.Sum512([]byte(new_password))
	decrypted, err := Decrypt(pw[:], wl.CipherKeys)
	if err != nil {
		return err
	}
	var pk PlainKeys
	err = json.Unmarshal([]byte(decrypted), &pk)
	if err != nil {
		return err
	}
	//Set keys
	for k := range pk.Keys {
		client.SetKeys(&Keys{OKey: []string{pk.Keys[k]}})
	}

	return nil
}

func SaveWalletFile(wallet_path string, wallet_filename string, password string, wallet_data *WalletData) error {
	// Validate data
	if password == "" {
		return errors.New("Password is not empty.")
	}
	if len(password) < 8 {
		return errors.New("Password length >= 8 character.")
	}
	if wallet_data == nil || wallet_data.Name == "" || wallet_data.PrivateKey == "" || wallet_data.PublicKey == "" {
		return errors.New("WalletData is invalid.")
	}

	// Serialize in memory, then save to disk
	//
	// This approach lessens the risk of a partially written wallet
	// if exceptions are thrown in serialization
	keys := make(map[string]string)
	keys[wallet_data.PublicKey] = wallet_data.PrivateKey
	salt, err := RandStringBytes(16)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}
	new_password := password + salt
	checksum := sha512.Sum512([]byte(new_password))

	var plainKeys PlainKeys
	plainKeys.Keys = keys
	copy(plainKeys.Checksum[:], checksum[:])
	plainData, err := json.Marshal(plainKeys)
	if err != nil {
		return err
	}
	plainTxt := string(plainData)
	cipherKeys, err := Encrypt(plainKeys.Checksum[:], plainTxt)

	if wallet_filename == "" {
		wallet_filename = wallet_data.Name + "-" + WalletName_
	}
	file_path := wallet_filename
	if wallet_path != "" {
		file_path = wallet_path + string(os.PathSeparator) + wallet_filename
	}

	wl := Wallet{CipherKeys: cipherKeys, CipherType: "aes-256-cbc", Name: wallet_data.Name, Salt: salt}
	data, err := json.Marshal(wl)
	if err == nil {
		f, _ := os.OpenFile(file_path, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
		defer f.Close()
		f.Write(data)
	}
	return err
}

func EncodeWallet(password string, wallet_data *WalletData) (string, error) {
	// Validate data
	if password == "" {
		return "", errors.New("Password is not empty.")
	}
	if len(password) < 8 {
		return "", errors.New("Password length >= 8 character.")
	}
	if wallet_data == nil || wallet_data.Name == "" || wallet_data.PrivateKey == "" || wallet_data.PublicKey == "" {
		return "", errors.New("WalletData is invalid.")
	}

	// Serialize in memory, then save to disk
	//
	// This approach lessens the risk of a partially written wallet
	// if exceptions are thrown in serialization
	keys := make(map[string]string)
	keys[wallet_data.PublicKey] = wallet_data.PrivateKey
	salt, err := RandStringBytes(16)
	if err != nil {
		// Serve an appropriately vague error to the
		// user, but log the details internally.
		panic(err)
	}
	new_password := password + salt
	checksum := sha512.Sum512([]byte(new_password))

	var plainKeys PlainKeys
	plainKeys.Keys = keys
	copy(plainKeys.Checksum[:], checksum[:])
	plainData, err := json.Marshal(plainKeys)
	if err != nil {
		return "", err
	}
	plainTxt := string(plainData)
	cipherKeys, err := Encrypt(plainKeys.Checksum[:], plainTxt)

	wl := Wallet{CipherKeys: cipherKeys, CipherType: "aes-256-cbc", Name: wallet_data.Name, Salt: salt}
	data, err := json.Marshal(wl)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func (client *Client) SetKeysFromEncodeWallet(wallet_json string, password string) error {
	if wallet_json == "" {
		return errors.New("Wallet json is not empty.")
	}
	if password == "" {
		return errors.New("Password is not empty.")
	}

	var wl *Wallet
	err := json.Unmarshal([]byte(wallet_json), &wl)
	if wl == nil || err != nil {
		return errors.New("Can not decode json wallet data.")
	}

	new_password := password + wl.Salt
	pw := sha512.Sum512([]byte(new_password))
	decrypted, err := Decrypt(pw[:], wl.CipherKeys)
	if err != nil {
		return err
	}
	var pk PlainKeys
	err = json.Unmarshal([]byte(decrypted), &pk)
	if err != nil {
		return err
	}
	//Set keys
	for k := range pk.Keys {
		client.SetKeys(&Keys{OKey: []string{pk.Keys[k]}})
	}

	return nil
}
