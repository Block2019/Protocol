package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type TxOutput struct { //Transaction output
	Value  int    //Value of Transaction
	PubKey string //Wallet Address
}

type TxInput struct { //References to previous outputs
	ID        []byte //Hash Transaction ID
	Out       int    // index which the previous output appears
	Signature string // Script that provides the data that is sued in the PubKey
}

type Transaction struct {
	ID      []byte     //Hash transaction ID
	Inputs  []TxInput  //Values entered into transaction
	Outputs []TxOutput //Values shown by the transaction
}

//Coinbase Function
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to '%s'", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

//SetID function
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//Function to cehck if the transaction is valid
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

//Function to unlock the data
func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data //Signature is the data that is used in the PubKey
}

//Function to unlock the data
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data //PubKey is the data that is used in the PubKey
}

//New Transaction function
func NewTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Println("Not enough funds")
		return nil
	}

	//Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from}) //Change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
