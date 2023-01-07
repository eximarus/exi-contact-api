package dy

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CursorSerialize(cursor map[string]types.AttributeValue) (*string, error) {
	if len(cursor) < 1 {
		return nil, nil
	}

	var d map[string]interface{}
	attributevalue.UnmarshalMap(cursor, &d)

	marshalled, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	cipherText, err := encryptCursor(marshalled)
	if err != nil {
		return nil, err
	}

	return aws.String(base64.StdEncoding.EncodeToString(cipherText)), nil
}

func CursorDeserialize(cursor *string) (map[string]types.AttributeValue, error) {
	if cursor == nil {
		return nil, nil
	}

	cipherText, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := decryptCursor(cipherText)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	err = json.Unmarshal(jsonBytes, &d)
	if err != nil {
		return nil, err
	}

	marshalled, err := attributevalue.MarshalMap(d)
	if err != nil {
		return nil, err
	}
	return marshalled, nil
}

func encryptCursor(cursor []byte) ([]byte, error) {
	cursor = pkcs5Padding(cursor, aes.BlockSize, len(cursor))

	key := os.Getenv("AES_ENCRYPTION_SECRET")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(cursor))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], cursor)

	return ciphertext, nil
}

func decryptCursor(ciphertext []byte) ([]byte, error) {
	key := os.Getenv("AES_ENCRYPTION_SECRET")

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err

	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	return pkcs5Trimming(ciphertext), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
