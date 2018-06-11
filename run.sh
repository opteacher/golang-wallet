export GOPATH=`pwd`
read -p "Whether install/update modules(Y/n default n): " input
if [ $input = "Y" ]
then /home/george/applications/go/bin/go get -u github.com/go-sql-driver/mysql
fi
read -p "Choise build type(main/test): " input
/home/george/applications/go/bin/go run $input.go