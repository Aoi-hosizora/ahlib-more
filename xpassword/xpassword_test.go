package xpassword

import (
	"log"
	"testing"
)

func TestPassword(t *testing.T) {
	password := []byte("123")

	encrypted, err := Encrypt(password, MinCost)
	log.Println(string(encrypted), err)
	check, err := Check(password, encrypted)
	log.Println(check, err)

	encrypted, err = Encrypt(password, DefaultCost)
	log.Println(string(encrypted), err)
	check, err = Check(password, encrypted)
	log.Println(check, err)

	encrypted, err = EncryptWithDefaultCost(password)
	log.Println(string(encrypted), err)
	check, err = Check(password, encrypted)
	log.Println(check, err)
}
