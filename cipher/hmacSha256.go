package cipher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CheckHmacSha256(message []byte, messageMAC string, key []byte) bool {
	return messageMAC == HmacSha256(message, key)
}

func HmacSha256(message []byte, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}
