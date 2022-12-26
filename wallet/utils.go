package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

/*
	Base58 is an algorithm invented with Bitcoin. Is derived from the base64 algorithm.
	It uses 6 less characters inside of its alphabet (0 O l I + /)- This 6 characters are
	easily confuse between them and since the public id are sensitive and we don't want users
	sending tokens to the wrong address
*/
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	if err != nil {
		log.Panic(err)
	}

	return decode
}
