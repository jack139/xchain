package main

import (
	"encoding/base64"
	crypto_rand "crypto/rand"
	"fmt"
	"golang.org/x/crypto/nacl/box"
	"io"
)

func main() {
	var senderPublicKey, senderPrivateKey, recipientPublicKey, recipientPrivateKey [32]byte

	// user1
	_senderPublicKey, _  := base64.StdEncoding.DecodeString("qyBsXnVKKjvFNxHBRudc3tCp8t8ymqBSF1Ga8qlfqFs=")
	_senderPrivateKey, _ := base64.StdEncoding.DecodeString("tgNfUoYkh9xKs1hVKs+5uXNetCxvDRRHBNmLMs5/NKk=")

	copy(senderPublicKey[:], _senderPublicKey[:32])
	copy(senderPrivateKey[:], _senderPrivateKey[:32])

	// user2
	_recipientPublicKey, _  := base64.StdEncoding.DecodeString("j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg=")
	_recipientPrivateKey, _ := base64.StdEncoding.DecodeString("gix62BaGqjD3ktH8G3lHdIuTTku6o2fHkkrJ5kLbID4=")

	copy(recipientPublicKey[:], _recipientPublicKey[:32])
	copy(recipientPrivateKey[:], _recipientPrivateKey[:32])

	// 加密
	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, &recipientPublicKey, &senderPrivateKey)

	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	msg := []byte("xxx")

	fmt.Printf("data=>%v,nonce=>%v,sharedEncryptKey=>%v\n", msg, nonce, *sharedEncryptKey)
	encrypted := box.SealAfterPrecomputation(nonce[:], msg, &nonce, sharedEncryptKey)

	// 解密
	var sharedDecryptKey [32]byte
	box.Precompute(&sharedDecryptKey, &senderPublicKey, &recipientPrivateKey)

	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	fmt.Printf("data=>%v,decryptNonce=>%v,decryptKey=>%v\n", encrypted[24:], decryptNonce, sharedDecryptKey)
	decrypted, ok := box.OpenAfterPrecomputation(nil, encrypted[24:], &decryptNonce, &sharedDecryptKey)
	if !ok {
		panic("decryption error")
	}
	fmt.Println(string(decrypted))
}