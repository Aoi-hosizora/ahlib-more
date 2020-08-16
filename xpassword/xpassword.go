package xpassword

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // The minimum allowable cost.
	MaxCost     int = 31 // The maximum allowable cost.
	DefaultCost int = 10 // The cost that will actually be set if a cost is below MinCost.
)

// Use bcrypt with cost to encrypt password.
// If the cost given is less than MinCost, the cost will be set to DefaultCost instead.
func Encrypt(password []byte, cost int) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// Use bcrypt with DefaultCost to encrypt password.
func EncryptWithDefaultCost(password []byte) ([]byte, error) {
	return Encrypt(password, DefaultCost)
}

// Check the password is the same.
func Check(password, encrypted []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(encrypted, password)
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
