package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const dbFile = "./tmp/blocks"

//BlockChain structure
type BlockChain struct {
	LastHash []byte //Last Hash of the BlockChain
	Database *badger.DB //Database

	//Old Implementation before using DB
	//Blocks []*Block //slice of blocks
	//The BlockChain is a slice of Blocks of type Block as shown above
}

//Initializing the BlockChain
func InitBlockChain() *BlockChain {
	var lastHash []byte
	opts := badger.DefaultOptions
	opts.Dir = dbFile
	opts.ValueDir = dbFile

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

	if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
		genesis := Genesis()
		fmt.Println("Genesis Block: ", genesis)
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err

	}
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value()
		return err
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain



	//return &BlockChain{[]*Block{Genesis()}} //Creating a new BlockChain with Genesis Block
}

//Handle the error
func Handle(err error) {
	if err != nil {
		panic(err)
	}
}


//Creating a New Block with its Data and Previous Hash
func NewBlock(data string, prevHash []byte) *Block {
// 	block := &Block{}         //Creating a new Block of the data structure
// 	block.Data = []byte(data) //Setting the Data
// 	block.PrevHash = prevHash //Setting the Previous Hash
// //	block.SetHash()           //Setting the Hash
// //We shall replace the above with the proof of work
// 	pow := NewProof(block) //Creating a new Proof of Work
// 	nonce, hash := pow.Run() //Running the Proof of Work
// 	block.Hash = hash[:] 	//Setting the Hash
// 	block.Nonce = nonce 	//Setting the Nonce
		
// 	return block
var lastHash []byte

err := db.View(func(txn *badger.Txn) error {
	item, err := txn.Get([]byte("lh"))
	Handle(err)
	lastHash, err = item.Value()
	return err
})
Handle(err)

	newBlock := NewBlock(data, lastHash)

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		lastHash = newBlock.Hash
		return err
	})

	Handle(err)

}

