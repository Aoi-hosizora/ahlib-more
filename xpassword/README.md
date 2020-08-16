# xpassword

### Functions

+ `Encrypt(password []byte, cost int) ([]byte, error)`
+ `EncryptWithDefaultCost(password []byte) ([]byte, error)`
+ `Check(password, encrypted []byte) (bool, error)`
