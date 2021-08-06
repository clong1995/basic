package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesEncrypt 加密
func AesEncrypt(orig, key []byte) ([]byte, error) {
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	orig = pKCS7Padding(orig, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	// 创建数组
	crypt := make([]byte, len(orig))
	// 加密
	blockMode.CryptBlocks(crypt, orig)

	return crypt, nil
}

// AesDecrypt 解密
func AesDecrypt(crypt, key []byte) ([]byte, error) {
	// 分组秘钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// 创建数组
	orig := make([]byte, len(crypt))
	// 解密
	blockMode.CryptBlocks(orig, crypt)
	// 去补全码
	orig = pKCS7UnPadding(orig)
	return orig, nil
}

//补码
func pKCS7Padding(orig []byte, blockSize int) []byte {
	padding := blockSize - len(orig)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(orig, padText...)
}

//去码
func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
