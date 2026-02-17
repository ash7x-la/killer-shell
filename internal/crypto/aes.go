package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func EncryptAES(plaintext []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...)
}

func DecryptAES(data []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	if len(data) < gcm.NonceSize() {
		return nil
	}
	nonce, ct := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	res, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil
	}
	return res
}
