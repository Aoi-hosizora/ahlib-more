package xpassword

import (
	"log"
	"testing"
)

func TestPassword(t *testing.T) {
	password := []byte("123")

	encrypted, err := EncryptPassword(password, MinCost)
	log.Println(string(encrypted), err)
	check, err := CheckPassword(password, encrypted)
	log.Println(check, err)

	encrypted, err = EncryptPassword(password, DefaultCost)
	log.Println(string(encrypted), err)
	check, err = CheckPassword(password, encrypted)
	log.Println(check, err)
}
