package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:])) //we get a []byte back 
	if err != nil {
		log.Panic(err)
	}
	return decode 	//No need to convert it to []byte
}

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input) // we get a string back

	return []byte(encode) //we need to convert it to []byte
}

//BASE58 algorith was invented by Bitcoin , 0OlI were removed to avoid people sending to wrong address


