package blockchain

import (
	"bytes"
	"crypto/sha256"
)

//Block structure
type Block struct {

	//Block header
	Hash     []byte 		//Hash of the Block
	PrevHash []byte 	//Previous Hash
	Data     []byte //Data of the Block
	Nonce    int //Proof of Work
}

//BlockChain structure
type BlockChain struct {
	Blocks []*Block //slice of blocks
	//The BlockChain is a slice of Blocks of type Block as shown above
}

//Creating the Hash
func (b *Block) SetHash() {
	info := bytes.Join([][]byte{b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

//Adding the Previous Hash to the Block
func (b *Block) SetPrevHash(prevHash []byte) {
	b.PrevHash = prevHash
}

//Adding the Data to the Block
func (b *Block) SetData(data []byte) {
	b.Data = data
}

//Get the Hash of the Block
func (b *Block) GetHash() []byte {
	return b.Hash
}

//Get the Previous Hash of the Block
func (b *Block) GetPrevHash() []byte {
	return b.PrevHash
}

//Get the Data of the Block
func (b *Block) GetData() []byte {
	return b.Data
}

// //Creating the BlockChain - a slice of Blocks called youthchain
// func (youthchain *BlockChain) GetBlocks() []*Block { // The blocks are of type Block
// 	return youthchain.blocks 	//returning the slice of blocks in a chain called youthchain
// }

//Creating a New Block with its Data and Previous Hash
func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{}         //Creating a new Block of the data structure
	block.Data = []byte(data) //Setting the Data
	block.PrevHash = prevHash //Setting the Previous Hash
//	block.SetHash()           //Setting the Hash
//We shall replace the above with the proof of work
	pow := NewProof(block) //Creating a new Proof of Work
	nonce, hash := pow.Run() //Running the Proof of Work
	block.Hash = hash[:] 	//Setting the Hash
	block.Nonce = nonce 	//Setting the Nonce
		
	return block
}



//The Genesis Block - Has no previous hash
func Genesis() *Block {
	return NewBlock("Genesis", []byte{}) //Empty Previous Hash
}

//Initializing the BlockChain
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}} //Creating a new BlockChain with Genesis Block
}

func (youthchain *BlockChain) AddBlock(data string) {
	prevBlock := youthchain.Blocks[len(youthchain.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.GetHash())
	youthchain.Blocks = append(youthchain.Blocks, newBlock)
}
