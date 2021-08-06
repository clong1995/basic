package cipher

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
)

// Base64EncryptBytes 加密byte[]为string
func Base64EncryptBytes(buf []byte) string {
	return base64.RawURLEncoding.EncodeToString(buf)
}

// Base64DecryptBytes 解密string为byte[]
func Base64DecryptBytes(base64Str string) ([]byte, error) {
	if base64Str == "" {
		err := fmt.Errorf("base64Str is empty")
		log.Println(err)
		return nil, err
	}
	return base64.RawURLEncoding.DecodeString(base64Str)
}

// Base64EncryptInt64 加密int64为string
func Base64EncryptInt64(number int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(number))
	return Base64EncryptBytes(b)
}

// Base64DecryptBytesInt 解密string为int64
func Base64DecryptBytesInt(base64Str string) int64 {
	buf, err := Base64DecryptBytes(base64Str)
	if err != nil {
		log.Println(err)
		return 0
	}

	if len(buf) < 7 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(buf))
}
