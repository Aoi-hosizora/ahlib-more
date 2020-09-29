package xjwt

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	xtesting.Equal(t, DefaultValidationError.Error(), "token is invalid")
	xtesting.Equal(t, DefaultValidationError.Errors, jwt.ValidationErrorClaimsInvalid)

	xtesting.True(t, CheckFlagError(DefaultValidationError, jwt.ValidationErrorClaimsInvalid))
	xtesting.False(t, CheckFlagError(DefaultValidationError, jwt.ValidationErrorExpired))

	xtesting.True(t, TokenExpired(jwt.NewValidationError("", jwt.ValidationErrorExpired)))
	xtesting.False(t, TokenExpired(DefaultValidationError))

	xtesting.True(t, TokenNotIssued(jwt.NewValidationError("", jwt.ValidationErrorIssuedAt)))
	xtesting.False(t, TokenNotIssued(DefaultValidationError))

	xtesting.True(t, TokenIssuerInvalid(jwt.NewValidationError("", jwt.ValidationErrorIssuer)))
	xtesting.False(t, TokenIssuerInvalid(DefaultValidationError))

	xtesting.True(t, TokenNotValidYet(jwt.NewValidationError("", jwt.ValidationErrorNotValidYet)))
	xtesting.False(t, TokenNotValidYet(DefaultValidationError))

	xtesting.True(t, TokenClaimsInvalid(DefaultValidationError))
	xtesting.False(t, TokenClaimsInvalid(jwt.NewValidationError("", jwt.ValidationErrorNotValidYet)))
}

func TestGenerateTokenAndParseToken(t *testing.T) {
	secret := []byte("A!B@C#D$E%F^G&")
	type userClaims struct {
		Uid uint64
		jwt.StandardClaims
	}

	now := time.Now()
	token, err := GenerateToken(&userClaims{
		Uid: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(2 * time.Second).Unix(),
			IssuedAt:  now.Add(1 * time.Second).Unix(),
		},
	}, secret)
	xtesting.Nil(t, err)

	// not issued
	_, err = ParseToken(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.True(t, TokenNotIssued(err))
	xtesting.False(t, TokenExpired(err))

	// issued
	time.Sleep(1100 * time.Millisecond)
	claims, err := ParseToken(token, secret, &userClaims{})
	cl, ok := claims.(*userClaims)
	xtesting.True(t, ok)
	xtesting.Equal(t, cl.ExpiresAt-now.Unix(), int64(2))
	xtesting.Equal(t, cl.IssuedAt-now.Unix(), int64(1))

	// expired
	time.Sleep(2000 * time.Millisecond)
	_, err = ParseToken(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.False(t, TokenNotIssued(err))
	xtesting.True(t, TokenExpired(err))
}
