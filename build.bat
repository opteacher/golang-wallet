set GOPATH=%CD%
set /p input=Whether install/update modules(Y/n):
if %input%==Y (
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/stretchr/testify
)
go build wallet.go