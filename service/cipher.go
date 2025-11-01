package service

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type Cipher interface {
	Encrypt(input []byte) ([]byte, error)
	Decrypt(input []byte) ([]byte, error)
}

type NoopCipher struct{}

func (noop NoopCipher) Encrypt(input []byte) ([]byte, error) {
	return input, nil
}

func (noop NoopCipher) Decrypt(input []byte) ([]byte, error) {
	return input, nil
}

type AesEcbCipher struct {
	cipherBlok cipher.Block
}

func NewAesEcbCipher(key []byte) AesEcbCipher {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return AesEcbCipher{
		cipherBlok: block,
	}
}

func (aesEcbCipher AesEcbCipher) Encrypt(input []byte) ([]byte, error) {
	if len(input)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("input length must be a multiple of the AES block size (16 bytes)")
	}

	ciphertext := make([]byte, len(input))
	for start := 0; start < len(input); start += aes.BlockSize {
		aesEcbCipher.cipherBlok.Encrypt(ciphertext[start:start+aes.BlockSize], input[start:start+aes.BlockSize])
	}
	return ciphertext, nil
}

func (aesEcbCipher AesEcbCipher) Decrypt(input []byte) ([]byte, error) {
	if len(input)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("input length must be a multiple of the AES block size (16 bytes)")
	}

	plaintext := make([]byte, len(input))
	for start := 0; start < len(input); start += aes.BlockSize {
		aesEcbCipher.cipherBlok.Decrypt(plaintext[start:start+aes.BlockSize], input[start:start+aes.BlockSize])
	}
	return plaintext, nil
}
