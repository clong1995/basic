package cipher

import (
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	encrypt, err := AesEncrypt([]byte("123"), []byte("crypto aes key !"))
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(Base64EncryptBytes(encrypt))
}
