# xjwt

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/golang-jwt/jwt

## Documents

### Types

+ None

### Variables

+ None

### Constants

+ None

### Functions

+ `func GenerateToken(method jwt.SigningMethod, claims jwt.Claims, key interface{}) (string, error)`
+ `func GenerateTokenWithHS256(claims jwt.Claims, secret []byte) (string, error)`
+ `func GenerateTokenWithHS384(claims jwt.Claims, secret []byte) (string, error)`
+ `func GenerateTokenWithHS512(claims jwt.Claims, secret []byte) (string, error)`
+ `func ParseToken(signedToken string, secret []byte, claims jwt.Claims) (*jwt.Token, error)`
+ `func ParseTokenClaims(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error)`
+ `func CheckValidationError(err error, flag uint32) bool`
+ `func IsAudienceError(err error) bool`
+ `func IsExpiredError(err error) bool`
+ `func IsIdError(err error) bool`
+ `func IsIssuedAtError(err error) bool`
+ `func IsIssuerError(err error) bool`
+ `func IsNotValidYetError(err error) bool`
+ `func IsTokenInvalidError(err error) bool`
+ `func IsClaimsInvalidError(err error) bool`

### Methods

+ None
