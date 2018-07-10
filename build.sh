export GOPATH=`pwd`
read -p "Whether install/update modules(Y/n default n): " input
if [ $input = "Y" ]
then \
/home/george/applications/go/bin/go get -u github.com/go-sql-driver/mysql |\
/home/george/applications/go/bin/go get -u github.com/stretchr/testify |\
/home/george/applications/go/bin/go get -u github.com/go-redis/redis
fi
/home/george/applications/go/bin/go build wallet.go