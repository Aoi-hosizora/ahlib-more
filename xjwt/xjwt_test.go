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

	secret := []byte("A!B@C#D$E%F^G&")
	token, err := GenerateToken(&userClaims{
		Uid: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Second).Unix(),
		},
	}, secret)
	log.Println(token, err)

	claims, err := ParseToken(token, secret, &userClaims{})
	log.Println(claims, err)
	c, ok := claims.(*userClaims)
	log.Println(c, ok)

	time.Sleep(3 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(IsTokenExpireError(err))

	time.Sleep(3 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(IsTokenExpireError(err))

	time.Sleep(1 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	log.Println(IsTokenExpireError(err))
}
