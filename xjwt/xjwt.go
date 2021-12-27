package xjwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken generates jwt token using given jwt.Claims, secret and jwt.SigningMethod.
func GenerateToken(method jwt.SigningMethod, claims jwt.Claims, key interface{}) (string, error) {
	tokenObj := jwt.NewWithClaims(method, claims)
	token, err := tokenObj.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GenerateTokenWithHS256 generates token using given jwt.Claims, secret and HS256 (HMAC SHA256, jwt.SigningMethodHS256) signing method.
func GenerateTokenWithHS256(claims jwt.Claims, secret []byte) (string, error) {
	return GenerateToken(jwt.SigningMethodHS256, claims, secret)
}

// GenerateTokenWithHS384 generates token using given jwt.Claims, secret and HS384 (HMAC SHA384, jwt.SigningMethodHS384) signing method.
func GenerateTokenWithHS384(claims jwt.Claims, secret []byte) (string, error) {
	return GenerateToken(jwt.SigningMethodHS384, claims, secret)
}

// GenerateTokenWithHS512 generates token using given jwt.Claims, secret and HS512 (HMAC SHA512, jwt.SigningMethodHS512) signing method.
func GenerateTokenWithHS512(claims jwt.Claims, secret []byte) (string, error) {
	return GenerateToken(jwt.SigningMethodHS512, claims, secret)
}

// ParseToken parses jwt token string to jwt.Token using given jwt.Claims and secret.
func ParseToken(signedToken string, secret []byte, claims jwt.Claims) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}
	tokenObj, err := jwt.ParseWithClaims(signedToken, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	return tokenObj, nil
}

// ParseTokenClaims parses jwt token string to jwt.Claims using given jwt.Claims and secret.
func ParseTokenClaims(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error) {
	tokenObj, err := ParseToken(signedToken, secret, claims)
	if err != nil {
		return nil, err
	}
	return tokenObj.Claims, nil
}

// CheckValidationError returns true if the given error is jwt.ValidationError with given flag.
func CheckValidationError(err error, flag uint32) bool {
	ve, ok := err.(*jwt.ValidationError)
	return ok && (ve.Errors&flag != 0)
}

// IsAudienceError checks error is an AUD (Audience) validation error.
func IsAudienceError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorAudience)
}

// IsExpiredError checks error is an EXP (Expires at) validation error.
func IsExpiredError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorExpired)
}

// IsIdError checks error is a JTI (Id) validation error.
func IsIdError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorId)
}

// IsIssuedAtError checks error is an IAT (Issued at) validation error.
func IsIssuedAtError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorIssuedAt)
}

// IsIssuerError checks error is an ISS (Issuer) validation error.
func IsIssuerError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorIssuer)
}

// IsNotValidYetError checks error is a NBF (Not before) validation error.
func IsNotValidYetError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorNotValidYet)
}

// // IsSubjectError checks error is a SUB (Subject) validation error.
// func IsSubjectError(err error) bool {
// 	return CheckValidationError(err, jwt.ValidationErrorSubject) // not found
// }

// IsTokenInvalidError checks error is an invalid token (could not be parsed) error.
func IsTokenInvalidError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorMalformed|jwt.ValidationErrorUnverifiable|jwt.ValidationErrorSignatureInvalid)
}

// IsClaimsInvalidError checks error is a generic claims validation error.
func IsClaimsInvalidError(err error) bool {
	return CheckValidationError(err, jwt.ValidationErrorClaimsInvalid)
}
