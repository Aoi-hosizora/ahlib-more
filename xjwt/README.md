# xjwt

### Functions

+ `var DefaultValidatorError jwt.ValidationError`
+ `GenerateToken(claims jwt.Claims, secret []byte) (string, error)`
+ `ParseToken(signedToken string, secret []byte, claims jwt.Claims) (jwt.Claims, error)`
+ `IsTokenExpireError(err error) bool`
