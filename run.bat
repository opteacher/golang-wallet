set GOPATH=%CD%
set /p input=Whether install/update modules(Y/n default n): 
if %input%=='Y' (
	go get -u github.com/go-sql-driver/mysql
)
set /p input=Choice build type(main/test):
go run %input%.go