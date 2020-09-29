package xpassword

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"log"
	"testing"
)

func TestPassword(t *testing.T) {
	password := []byte("123")

	for _, cost := range []int{1, MinCost, 7, DefaultCost, 13} { // 1 4 7 10 13
		encrypted, err := Encrypt(password, cost)
		xtesting.Nil(t, err)
		log.Println(cost, ":", string(encrypted))
		check, err := Check(password, encrypted) // true nil
		xtesting.Nil(t, err)
		xtesting.True(t, check)
	}

	encrypted, err := Encrypt(password, MaxCost+100)
	xtesting.NotNil(t, err)

	encrypted, err = EncryptWithDefaultCost(password)
	xtesting.Nil(t, err)
	log.Println(DefaultCost, ":", string(encrypted))
	check, err := Check(password, encrypted) // true nil
	xtesting.Nil(t, err)
	xtesting.True(t, check)

	check, err = Check(password, []byte{}) // false, err
	xtesting.NotNil(t, err)
	xtesting.False(t, check)

	check, err = Check([]byte("miss"), encrypted) // false, nil
	xtesting.Nil(t, err)
	xtesting.False(t, check)
}
