module github.com/Aoi-hosizora/ahlib-more

go 1.15

require (
	github.com/Aoi-hosizora/ahlib v0.0.0-00010101000000-000000000000
	github.com/ah-forklib/strftime v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/text v0.3.0
)

replace (
	github.com/Aoi-hosizora/ahlib => ../ahlib
	github.com/ah-forklib/strftime => ../_ref/strftime
)
