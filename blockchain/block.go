package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

//Block structure
type Block struct {

	//Block header
	Hash     []byte //Hash of the Block
	PrevHash []byte //Previous Hash
	//Data     []byte //Data of the Block - This was replaced woth Transaction Data
	Transactions []*Transaction //Transaction of the Block
	Nonce    int    //Proof of Work

}

//Hash the Block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// //Creating the Hash
// func (b *Block) SetHash() {
// 	info := bytes.Join([][]byte{b.PrevHash, b.Transactions}, []byte{})
// 	hash := sha256.Sum256(info)
// 	b.Hash = hash[:]
// }

//Adding the Previous Hash to the Block
// func (b *Block) SetPrevHash(prevHash []byte) {
// 	b.PrevHash = prevHash
// }

//Adding the Data to the Block
// func (b *Block) SetData(data []byte) {
// 	b.Transaction = data
// }

//Get the Hash of the Block
// func (b *Block) GetHash() []byte {
// 	return b.Hash
// }

//Get the Previous Hash of the Block
// func (b *Block) GetPrevHash() []byte {
// 	return b.PrevHash
// }

//Get the Data of the Block
// func (b *Block) GetData() []byte {
// 	return b.Transaction
// }

// //Creating the BlockChain - a slice of Blocks called youthchain
// func (youthchain *BlockChain) GetBlocks() []*Block { // The blocks are of type Block
// 	return youthchain.blocks 	//returning the slice of blocks in a chain called youthchain
// }

//The Genesis Block - Has no previous hash
func Genesis(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}) //Empty Previous Hash
}


//Serializing the Block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer            //Buffer is a type that implements the io.Writer interface
	encoder := gob.NewEncoder(&result) //encoder is a type that implements the io.Writer interface
	err := encoder.Encode(b)           //Encode the block
	Handle(err)                        //Handle Errors
	return result.Bytes()              //Return the bytes
}

//Deserializing the Block
func Deserialize(data []byte) *Block {
	var block Block                                  //block is a type that implements the Block interface
	decoder := gob.NewDecoder(bytes.NewReader(data)) //decoder is a type that implements the io.Reader interface
	err := decoder.Decode(&block)                    //Decode the block
	Handle(err)                                      //Handle Errors
	return &block                                    //Return the block
}


