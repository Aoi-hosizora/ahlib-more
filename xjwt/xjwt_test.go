package xjwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	type userClaims struct {
		Uid uint64
		jwt.StandardClaims
	}

	log.Println(DefaultValidationError.Error())
	log.Println(CheckFlagError(DefaultValidationError, jwt.ValidationErrorClaimsInvalid))

	secret := []byte("A!B@C#D$E%F^G&")
	token, err := GenerateToken(&userClaims{
		Uid: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
			IssuedAt:  time.Now().Add(1 * time.Second).Unix(),
		},
	}, secret)
	log.Println(token, err)

	claims, err := ParseToken(token, secret, &userClaims{})
	log.Println(TokenExpired(err), TokenNotIssued(err)) // false true
	c, ok := claims.(*userClaims)
	log.Println(c, ok) // nil false

	time.Sleep(3 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(TokenExpired(err), TokenNotIssued(err)) // false false

	time.Sleep(3 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(TokenExpired(err), TokenNotIssued(err)) // ? false

	time.Sleep(1 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(TokenExpired(err), TokenNotIssued(err)) // true false

	_ = TokenIssuerInvalid
	_ = TokenNotValidYet
}
