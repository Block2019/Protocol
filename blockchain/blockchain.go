package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

const dbFile = "./tmp/blocks/MANIFEST"
const dbPath = "./tmp/blocks"
const genesisCoinbaseData = "Amon has created a blockchain"

//BlockChain structure
type BlockChain struct {
	LastHash []byte     //Last Hash of the BlockChain
	Database *badger.DB //Database

	//Old Implementation before using DB
	//Blocks []*Block //slice of blocks
	//The BlockChain is a slice of Blocks of type Block as shown above
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

//BLockchain iterator
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

//Next function to be used in interator
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)

		var valCopy []byte
		valCopy, err = item.ValueCopy(valCopy)
		Handle(err)
		block = Deserialize(valCopy)
		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block

}

//Database check function
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true

}

//Initializing the BlockChain
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if dbExists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	// opts := badger.DefaultOptions
	// opts.Dir = dbPath
	// opts.ValueDir = dbPath

	// db, err := badger.Open(opts)

	// Open the Badger database located in the /tmp/badger directory.

	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(dbPath))

	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

		cbtx := CoinbaseTx(address, genesisCoinbaseData)
		genesis := Genesis(cbtx)
		fmt.Println("Genesis Block: ", genesis)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err

	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain //

	//return &BlockChain{[]*Block{Genesis()}} //Creating a new BlockChain with Genesis Block
}

//Continue the Blockchain Function
func ContinueBlockChain(address string) *BlockChain {
	if !dbExists() {
		fmt.Println("No existing blockchain found")
		runtime.Goexit()
	}

	var lastHash []byte

	// Open the Badger database located in the /tmp/badger directory.

	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions(dbPath))

	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		//lastHash, err = item.Value()
		lastHash, err = item.ValueCopy(lastHash)
		Handle(err)

		return err
	})

	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

//Check for unspent outputs
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlocked(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return unspentTXs
}

//Find unspent transaction
func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

//Find Spendable Outputs
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

//Handle the error
func Handle(err error) {
	if err != nil {
		panic(err)
	}
}

//Creating a New Block with its Data and Previous Hash
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, prevHash, txs, 0} //Creating a new Block of the data structure
	// block.Transaction = txs //Setting the Data
	// block.PrevHash = prevHash //Setting the Previous Hash
	//	block.SetHash()           //Setting the Hash

	//We shall replace the above with the proof of work
	pow := NewProof(block)   //Creating a new Proof of Work
	nonce, hash := pow.Run() //Running the Proof of Work
	block.Hash = hash[:]     //Setting the Hash
	block.Nonce = nonce      //Setting the Nonce

	return block

}

func (youthchain *BlockChain) AddBlock(transaction []*Transaction) {
	var lastHash []byte

	err := youthchain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(lastHash)
		Handle(err)
		return err
	})
	Handle(err)

	newBlock := NewBlock(transaction, lastHash)

	err = youthchain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		lastHash = newBlock.Hash
		return err
	})

	Handle(err)

}
