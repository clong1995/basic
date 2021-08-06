package cipher

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CheckPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(bytes.Join([][]byte{[]byte("$2a$10$"), hashedPassword}, []byte("")), password)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("id or password mistake")
		return false
	}
	return true
}

func Password(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("password mistake")
		return ""
	}
	return string(hash)[7:]
}
