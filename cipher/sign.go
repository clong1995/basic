package cipher

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

// Sign 签名
func Sign(message, ak []byte) string {
	var buffer bytes.Buffer
	buffer.Write(message)
	buffer.Write(ak)
	sum := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(sum[:])
}

// CheckSign 检查签名
func CheckSign(signature string, message, ak []byte) bool {
	return signature == Sign(message, ak)
}
