package xjwt

import (
	"github.com/dgrijalva/jwt-go"
)

var DefaultValidatorError = jwt.ValidationError{}

func GenerateToken(claims jwt.Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}
	token, err := jwt.ParseWithClaims(signedToken, claims, keyFunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, DefaultValidatorError
	}

	return token.Claims, nil
}

func IsTokenExpireError(err error) bool {
	if err == nil {
		return false
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}
