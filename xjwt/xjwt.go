package xjwt

import (
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates token with jwt.SigningMethodHS256.
func GenerateToken(claims jwt.Claims, secret []byte) (string, error) {
	method := jwt.SigningMethodHS256
	return GenerateTokenWithMethod(method, claims, secret)
}

// GenerateTokenWithMethod uses a jwt.SigningMethod to generate token.
func GenerateTokenWithMethod(method jwt.SigningMethod, claims jwt.Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses jwt token using given secret.
func ParseToken(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}
	token, err := jwt.ParseWithClaims(signedToken, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	return token.Claims, nil
}

// Default validation error, use jwt.ValidationErrorClaimsInvalid.
var DefaultValidationError = jwt.NewValidationError("token is invalid", jwt.ValidationErrorClaimsInvalid)

// CheckFlagError checks standard Claim validation errors.
func CheckFlagError(err error, flag uint32) bool {
	if err == nil {
		return false
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&flag != 0
	}
	return false
}

// EXP validation failed.
func TokenExpired(err error) bool {
	return CheckFlagError(err, jwt.ValidationErrorExpired)
}

// IAT validation failed.
func TokenNotIssued(err error) bool {
	return CheckFlagError(err, jwt.ValidationErrorIssuedAt)
}

// ISS validation failed.
func TokenIssuerInvalid(err error) bool {
	return CheckFlagError(err, jwt.ValidationErrorIssuer)
}

// NBF validation failed.
func TokenNotValidYet(err error) bool {
	return CheckFlagError(err, jwt.ValidationErrorNotValidYet)
}

// Generic claims validation error.
func TokenClaimsInvalid(err error) bool {
	return CheckFlagError(err, jwt.ValidationErrorClaimsInvalid)
}
