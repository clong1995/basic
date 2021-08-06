package cipher

import (
	"testing"
)

func TestCheckPassword(t *testing.T) {
	//flag := CheckPassword([]byte("HLb5p8MRmb5mxg3xSCdNwekLEf7LR/yZu7lvWniCOPD9OSmHv9McS"),[]byte("123456"))
	flag := CheckPassword([]byte("O2WjG0h15JaDdE2PAIxRiu6B6w/.W1LpyBJeuNg/0JRjvJleynAx6"), []byte("ABAATnRKhhM123456"))
	t.Log(flag)
}

func TestPassword(t *testing.T) {
	password := Password([]byte("123456")) //HLb5p8MRmb5mxg3xSCdNwekLEf7LR/yZu7lvWniCOPD9OSmHv9McS
	t.Log(password)
}
