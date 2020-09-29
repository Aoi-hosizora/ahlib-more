package xjwt

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	xtesting.Equal(t, DefaultValidationError.Error(), "token is invalid")
	xtesting.Equal(t, DefaultValidationError.Errors, jwt.ValidationErrorClaimsInvalid)

	xtesting.False(t, CheckFlagError(nil, jwt.ValidationErrorExpired))
	xtesting.False(t, CheckFlagError(fmt.Errorf("other error"), jwt.ValidationErrorExpired))
	xtesting.False(t, CheckFlagError(DefaultValidationError, jwt.ValidationErrorExpired))
	xtesting.True(t, CheckFlagError(DefaultValidationError, jwt.ValidationErrorClaimsInvalid))

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

type fakeMethod struct{}

func (f fakeMethod) Verify(string, string, interface{}) error {
	return nil
}

func (f fakeMethod) Sign(string, interface{}) (string, error) {
	return "", fmt.Errorf("fake error")
}

func (f fakeMethod) Alg() string {
	return ""
}

func TestGenerateToken(t *testing.T) {
	fake := &fakeMethod{}
	_, err := GenerateTokenWithMethod(fake, &jwt.StandardClaims{}, []byte{})
	xtesting.NotNil(t, err)

	_, err = GenerateToken(&jwt.StandardClaims{}, []byte{})
	xtesting.Nil(t, err)
}

func TestToken(t *testing.T) {
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
