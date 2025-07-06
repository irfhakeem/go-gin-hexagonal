package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/config"
	"net/url"
)

type AESEncryptor struct {
	key string
	iv  string
}

func NewAESEncryptor(cfg config.AESConfig) ports.Encryptor {
	return &AESEncryptor{
		key: cfg.Key,
		iv:  cfg.IV,
	}
}

func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	blockSize := aes.BlockSize
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plainTextBlock := append([]byte(plaintext), padtext...)

	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(e.iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	decodedCiphertext, err := url.QueryUnescape(ciphertext)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(decodedCiphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	if len(ciphertextBytes) == 0 {
		return "", fmt.Errorf("ciphertext is empty")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(e.iv))
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	length := len(ciphertextBytes)
	unpadding := int(ciphertextBytes[length-1])

	if unpadding > aes.BlockSize || unpadding == 0 || unpadding > length {
		return "", fmt.Errorf("invalid padding")
	}

	for i := length - unpadding; i < length; i++ {
		if ciphertextBytes[i] != byte(unpadding) {
			return "", fmt.Errorf("invalid padding")
		}
	}

	plaintext := ciphertextBytes[:(length - unpadding)]

	return string(plaintext), nil
}
