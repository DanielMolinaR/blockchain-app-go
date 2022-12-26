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
	version        = byte(0x00) // Hexadecimal representatin of 0
)

/*
	A wallet is made of 2 keys, public and private key. The private key is the identifier for each
	of the accounts inside of the blockchain. Each private key needs to be completely random and unique.
	The public key is the key that can be given to other users. It is the key used to derive the address,
	which is the address that we use to send and receive data in the blockchain.

	-----------------------------------------------------------------------------------------------------

	ECDSA (Elliptic Curve Digital Signature Algorithm) is a Digital Signature Algorithm (DSA) which uses
	keys derived from elliptic curve cryptography (ECC). It is a particularly efficient equation based on
	public key cryptography (PKC). Thanks to ECDSA we can generate up to 10^77 different keys
*/
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (wallet Wallet) Address() []byte {
	pubKeyHashed := PublicKeyHash(wallet.PublicKey)

	versionedHash := append([]byte{version}, pubKeyHashed...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	fmt.Printf("pub key: %x\n", wallet.PublicKey)
	fmt.Printf("pub hash: %x\n", pubKeyHashed)
	fmt.Printf("address: %s\n", address)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256() // Output = 256 bytes

	private, err := ecdsa.GenerateKey(curve, rand.Reader) // Generates the private key with the curve and a random number generator
	if err != nil {
		log.Panic(err)
	}

	/*
		For creating the public key we use the concept of the eliptic curve multiplication by  picking
		values in the eliptic curve at random and we take that X and Y values, we convert them into
		bytes and append them together
	*/
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub

}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()

	return &Wallet{private, public}

}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()          // Create a ripemd husher
	_, err := hasher.Write(pubHash[:]) // Write the public key hashed into the hasher
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil) // Hash in ripemd the public key already hashed in sha256

	return publicRipMD
}

func Checksum(pubKeyHashed []byte) []byte {
	pubHash := sha256.Sum256(pubKeyHashed)
	pubHash = sha256.Sum256(pubHash[:])

	return pubHash[:checksumLength]
}
