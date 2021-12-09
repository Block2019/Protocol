package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

//Take date from the block

//Create a counter which starts from 0

//Create a hash of the data plus the counter

//Check if the hash its meets the requirements

//Requirments
//The first few bytes must contain 0s

const Difficulty = 1

type ProofOfWork struct { //Proof of work structure
	Block  *Block   //Block to mine
	Target *big.Int //Target to meet	
}

func NewProof(b *Block) *ProofOfWork { //Create a new proof of work
	target := big.NewInt(1)                  //Create a new target
	target.Lsh(target, uint(256-Difficulty)) //Left shift the target

	pow := &ProofOfWork{
		Block:  b,
		Target: target,
	} //Create a new proof of work

	return pow //Return the proof of work
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte { //Prepare the data to be hashed
	data := bytes.Join( //Join the data together
		[][]byte{
			pow.Block.PrevHash,       //Previous hash
			pow.Block.Data,           //Data
			ToHex(int64(nonce)),      //Nonce
			ToHex(int64(Difficulty)), //Difficulty
		},
		[]byte{},
	)

	return data //Return the data
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)                        //Create a new buffer
	err := binary.Write(buff, binary.BigEndian, num) //Write the number to the buffer
	if err != nil {
		panic(err) //Panic if there is an error
	}

	return buff.Bytes() //Return the buffer
}

func Validate(b *Block, pow *ProofOfWork) bool { //Validate the proof of work
	var intHash big.Int //Create a new big int

	data := pow.PrepareData(b.Nonce) //Prepare the data
	hash := sha256.Sum256(data)      //Hash the data

	intHash.SetBytes(hash[:]) //Set the hash to a big int

	return intHash.Cmp(pow.Target) == -1 //Return if the hash is less than the target
}

func (pow *ProofOfWork) Run() (int, []byte) { //Run the proof of work
	var intHash big.Int //Create a new big int
	var hash [32]byte   //Create a new hash

	nonce := 0 //Set the nonce to 0

	for nonce < math.MaxInt64 { //Loop until the nonce is greater than the max int
		data := pow.PrepareData(nonce) //Prepare the data
		hash = sha256.Sum256(data)     //Hash the data

		fmt.Printf("\r%x", hash)  //Print the hash
		intHash.SetBytes(hash[:]) //Set the hash to a big int

		if intHash.Cmp(pow.Target) == -1 { //Check if the hash is less than the target
			break //Break the loop
		} else {
			nonce++ //Increment the nonce
		}

	}

	fmt.Println()         //Print a new line
	return nonce, hash[:] //Return the nonce and the hash
}

//Function to validate the proof of work
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int //Create a new big int

	data := pow.PrepareData(pow.Block.Nonce) //Prepare the data
	hash := sha256.Sum256(data)              //Hash the data
	intHash.SetBytes(hash[:])                //Set the hash to a big int

	return intHash.Cmp(pow.Target) == -1 //Return if the hash is less than the target
}


// func (pow *ProofOfWork) Validate() bool {
//     var intHash big.Int

//     data := pow.InitNonce(pow.Block.Nonce)

//     hash := sha256.Sum256(data)
//     intHash.SetBytes(hash[:])

//     return intHash.Cmp(pow.Target) == -1
// }
