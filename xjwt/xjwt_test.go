package xjwt

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	for _, tc := range []struct {
		giveMethod jwt.SigningMethod
		wantSecret interface{}
		wantError  bool
	}{
		{jwt.SigningMethodNone, []byte{}, true},
		{jwt.SigningMethodHS256, []byte{}, false},
		{jwt.SigningMethodHS384, []byte{}, false},
		{jwt.SigningMethodHS512, []byte{}, false},
		{jwt.SigningMethodES256, []byte{}, true},
		{jwt.SigningMethodES384, []byte{}, true},
		{jwt.SigningMethodES512, []byte{}, true},
	} {
		_, err := GenerateToken(tc.giveMethod, &jwt.RegisteredClaims{}, tc.wantSecret)
		if tc.wantError {
			xtesting.NotNil(t, err)
		} else {
			xtesting.Nil(t, err)
		}
	}

	for _, tc := range []struct {
		giveFn     func(jwt.Claims, []byte) (string, error)
		giveSecret []byte
		wantError  bool
	}{
		{GenerateTokenWithHS256, []byte{}, false},
		{GenerateTokenWithHS256, []byte{'#'}, false},
		{GenerateTokenWithHS384, []byte{}, false},
		{GenerateTokenWithHS384, []byte{'#'}, false},
		{GenerateTokenWithHS512, []byte{}, false},
		{GenerateTokenWithHS512, []byte{'#'}, false},
	} {
		token, err := tc.giveFn(&jwt.RegisteredClaims{}, tc.giveSecret)
		if tc.wantError {
			xtesting.NotNil(t, err)
		} else {
			xtesting.Nil(t, err)
			xtesting.NotBlankString(t, token)
		}
	}
}

func TestParseToken(t *testing.T) {
	secret := []byte("A!B@C#D$E%F^G&")
	type userClaims struct {
		Uid      uint64
		Username string
		jwt.RegisteredClaims
	}
	uid := uint64(20)
	username := "test user"
	now := time.Now()

	claims := &userClaims{
		Uid:      uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    username,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(time.Second)),
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Second)),
		},
	}
	// | now | +1s | +2s | +3s |
	// |-----|-----------|-----|
	//   NBF       OK      EXP
	token, err := GenerateTokenWithHS256(claims, secret)
	xtesting.Nil(t, err)

	// 1. NBF
	_, err = ParseToken(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.True(t, IsNotValidYetError(err))

	_, err = ParseTokenClaims(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.True(t, IsNotValidYetError(err))

	// 2. Valid
	time.Sleep(time.Second)
	parsedToken, err := ParseToken(token, secret, &userClaims{})
	xtesting.Nil(t, err)
	xtesting.Equal(t, parsedToken.Claims.(*userClaims).Uid, uid)
	xtesting.Equal(t, parsedToken.Claims.(*userClaims).Username, username)
	xtesting.Equal(t, parsedToken.Claims.(*userClaims).Issuer, username)
	xtesting.Equal(t, parsedToken.Claims.(*userClaims).IssuedAt.Unix(), now.Unix())

	parsedClaims, err := ParseTokenClaims(token, secret, &userClaims{})
	xtesting.Nil(t, err)
	xtesting.Equal(t, parsedClaims.(*userClaims).Uid, uid)
	xtesting.Equal(t, parsedClaims.(*userClaims).Username, username)
	xtesting.Equal(t, parsedClaims.(*userClaims).Issuer, username)
	xtesting.Equal(t, parsedClaims.(*userClaims).IssuedAt.Unix(), now.Unix())

	// 3. EXP
	time.Sleep(2 * time.Second)
	_, err = ParseToken(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.True(t, IsExpiredError(err))

	_, err = ParseTokenClaims(token, secret, &userClaims{})
	xtesting.NotNil(t, err)
	xtesting.True(t, IsExpiredError(err))
}

func TestValidationError(t *testing.T) {
	for _, tc := range []struct {
		giveFn  func(error) bool
		giveErr error
		want    bool
	}{
		{IsAudienceError, nil, false},
		{IsAudienceError, errors.New(""), false},
		{IsAudienceError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsAudienceError, jwt.NewValidationError("", jwt.ValidationErrorAudience), true},
		{IsAudienceError, jwt.NewValidationError("", jwt.ValidationErrorAudience|jwt.ValidationErrorMalformed), true},

		{IsExpiredError, nil, false},
		{IsExpiredError, errors.New(""), false},
		{IsExpiredError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsExpiredError, jwt.NewValidationError("", jwt.ValidationErrorExpired), true},
		{IsExpiredError, jwt.NewValidationError("", jwt.ValidationErrorExpired|jwt.ValidationErrorMalformed), true},

		{IsIdError, nil, false},
		{IsIdError, errors.New(""), false},
		{IsIdError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsIdError, jwt.NewValidationError("", jwt.ValidationErrorId), true},
		{IsIdError, jwt.NewValidationError("", jwt.ValidationErrorId|jwt.ValidationErrorMalformed), true},

		{IsIssuedAtError, nil, false},
		{IsIssuedAtError, errors.New(""), false},
		{IsIssuedAtError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsIssuedAtError, jwt.NewValidationError("", jwt.ValidationErrorIssuedAt), true},
		{IsIssuedAtError, jwt.NewValidationError("", jwt.ValidationErrorIssuedAt|jwt.ValidationErrorMalformed), true},

		{IsIssuerError, nil, false},
		{IsIssuerError, errors.New(""), false},
		{IsIssuerError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsIssuerError, jwt.NewValidationError("", jwt.ValidationErrorIssuer), true},
		{IsIssuerError, jwt.NewValidationError("", jwt.ValidationErrorIssuer|jwt.ValidationErrorMalformed), true},

		{IsNotValidYetError, nil, false},
		{IsNotValidYetError, errors.New(""), false},
		{IsNotValidYetError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsNotValidYetError, jwt.NewValidationError("", jwt.ValidationErrorNotValidYet), true},
		{IsNotValidYetError, jwt.NewValidationError("", jwt.ValidationErrorNotValidYet|jwt.ValidationErrorMalformed), true},

		{IsTokenInvalidError, nil, false},
		{IsTokenInvalidError, errors.New(""), false},
		{IsTokenInvalidError, jwt.NewValidationError("", jwt.ValidationErrorClaimsInvalid), false},
		{IsTokenInvalidError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), true},
		{IsTokenInvalidError, jwt.NewValidationError("", jwt.ValidationErrorUnverifiable), true},
		{IsTokenInvalidError, jwt.NewValidationError("", jwt.ValidationErrorSignatureInvalid), true},
		{IsTokenInvalidError, jwt.NewValidationError("", jwt.ValidationErrorMalformed|jwt.ValidationErrorClaimsInvalid), true},

		{IsClaimsInvalidError, nil, false},
		{IsClaimsInvalidError, errors.New(""), false},
		{IsClaimsInvalidError, jwt.NewValidationError("", jwt.ValidationErrorMalformed), false},
		{IsClaimsInvalidError, jwt.NewValidationError("", jwt.ValidationErrorClaimsInvalid), true},
		{IsClaimsInvalidError, jwt.NewValidationError("", jwt.ValidationErrorClaimsInvalid|jwt.ValidationErrorMalformed), true},
	} {
		xtesting.Equal(t, tc.giveFn(tc.giveErr), tc.want)
	}
}
