package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"
)

// GetHmacResult generate HMAC value
func GetHmacResult(key string, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	encoder := base64.StdEncoding
	h.Write([]byte(encoder.EncodeToString([]byte(data))))
	return encoder.EncodeToString(h.Sum(nil))
}

// Encrypt encrypt msg with the given key.
func Encrypt(key string, msg string) (encrypted string, err error) {
	block, err := aes.NewCipher([]byte(key[0:16]))
	if err != nil {
		return "", err
	}

	encoder := base64.StdEncoding

	ecb := cipher.NewCBCEncrypter(block, getIvParam(key))

	content := []byte(encoder.EncodeToString([]byte(msg)))

	rand.Seed(time.Now().UnixNano())
	rnd := strconv.Itoa(rand.Intn(89999999) + 10000000)
	salt := []byte(rnd)

	content = append(salt, content...)

	content = getPKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return encoder.EncodeToString(crypted), nil
}

// Decrypt decrypt msg with the given key.
func Decrypt(key string, crypt string) (decrypted string, err error) {
	block, err := aes.NewCipher([]byte(key[0:16]))
	if err != nil {
		return "", err
	}

	encoder := base64.StdEncoding

	b, err := encoder.DecodeString(crypt)
	if err != nil {
		return "", err
	}

	ecb := cipher.NewCBCDecrypter(block, getIvParam(key))
	dec := make([]byte, len(b))
	ecb.CryptBlocks(dec, b)
	decrypt := getPKCS5Trimming(dec)
	decrypt = decrypt[8:len(decrypt)]

	b, err = encoder.DecodeString(string(decrypt))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getPKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func getPKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func getIvParam(key string) []byte {
	var param []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	b := []byte(key)
	len := len(b)

	for i := 0; i < 16 && i < len; i++ {
		param[i] = b[i]
	}

	return param
}
