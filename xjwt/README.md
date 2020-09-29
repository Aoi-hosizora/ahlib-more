# xjwt

### Dependencies

+ github.com/dgrijalva/jwt-go

### Functions

+ `var DefaultValidatorError jwt.ValidationError`
+ `GenerateToken(claims jwt.Claims, secret []byte) (string, error)`
+ `GenerateTokenWithMethod(method jwt.SigningMethod, claims jwt.Claims, secret []byte) (string, error)`
+ `ParseToken(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error)`
+ `CheckFlagError(err error, flag uint32) bool`
+ `TokenExpired(err error) bool`
+ `TokenNotIssued(err error) bool`
+ `TokenIssuerInvalid(err error) bool`
+ `TokenNotValidYet(err error) bool`
