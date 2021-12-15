package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

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
