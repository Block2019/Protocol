package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "./tmp/wallet.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) SaveToFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	//gob.Register(map[string]interface{}{})
	//gob.Register(Wallet{})

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) LoadFromFile() {

	var wallets Wallets
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return
	}

	var content []byte
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&ws)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

}

func CreateWallet() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err

}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets
	var content []byte
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address

}
