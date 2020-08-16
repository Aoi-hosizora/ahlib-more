package xpassword

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // The minimum allowable cost as passed in to bcrypt.GenerateFromPassword.
	MaxCost     int = 31 // The maximum allowable cost as passed in to bcrypt.GenerateFromPassword.
	DefaultCost int = 10 // The cost that will actually be set if a cost below MinCost is passed into bcrypt.GenerateFromPassword.
)

func EncryptPassword(password []byte, cost int) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func CheckPassword(password, encrypted []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(encrypted, password)
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
