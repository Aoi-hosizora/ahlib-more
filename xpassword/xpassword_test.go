package xpassword

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"log"
	"testing"
)

func TestPassword(t *testing.T) {
	password := []byte("123")

	_ = MaxCost
	for _, cost := range []int{1, MinCost, 7, DefaultCost, 13} { // 1 4 7 10 13
		encrypted, err := Encrypt(password, cost)
		xtesting.Nil(t, err)
		log.Println(cost, ":", string(encrypted))
		check, err := Check(password, encrypted)
		xtesting.Nil(t, err)
		xtesting.True(t, check)
	}

	encrypted, err := EncryptWithDefaultCost(password)
	xtesting.Nil(t, err)
	log.Println(DefaultCost, ":", string(encrypted))
	check, err := Check(password, encrypted)
	xtesting.Nil(t, err)
	xtesting.True(t, check)
}
