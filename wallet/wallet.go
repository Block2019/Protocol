package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // Elliptic curve digital signing alorithm private key
	PublicKey  []byte
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {

	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func NewWallet() *Wallet {

	private, public := NewKeyPair()

	wallet := Wallet{private, public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {

	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()

	_, err := hasher.Write(pubHash[:])

	if err != nil {
		log.Panic(err)
	}

	publicRipemd := hasher.Sum(nil)

	return publicRipemd

}

func Checksum(payload []byte) []byte {

	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	return second[:checksumLength]
}

func (w Wallet) GetAddress() []byte {

	pubKeyHash := PublicKeyHash(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)

	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)

	fmt.Printf("address: %x\n", address)
	fmt.Printf("pub hash: %x\n", pubKeyHash)
	fmt.Printf("pub Key: %x\n", w.PublicKey)

	return address
}
