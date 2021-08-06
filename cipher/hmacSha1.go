package cipher

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func CheckHmacSha1(message []byte, messageMAC string, key []byte) bool {
	return messageMAC == HmacSha1(message, key)
}

func HmacSha1(message []byte, key []byte) string {
	expectedMAC := HmacSha1Byte(message, key)
	return hex.EncodeToString(expectedMAC)
}

func HmacSha1Byte(message []byte, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
