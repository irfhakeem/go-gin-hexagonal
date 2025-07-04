package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"go-gin-hexagonal/internal/domain/ports"
	"go-gin-hexagonal/pkg/config"
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

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
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
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(e.key))

	if err != nil {
		return "", err
	}

	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(e.iv))
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	length := len(ciphertextBytes)
	unpadding := int(ciphertextBytes[length-1])

	plaintext := ciphertextBytes[:(length - unpadding)]

	return string(plaintext), nil
}
