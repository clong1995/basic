package cipher

import (
	"testing"
)

func TestHmacSha256(t *testing.T) {
	encrypt := HmacSha256([]byte("message"), []byte("key"))
	t.Log(encrypt)
}
