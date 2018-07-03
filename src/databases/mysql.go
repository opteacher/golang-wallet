package databases

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"utils"
	"strings"
	"net/url"
	"time"
)

var PARAMS = []string {
	"charset=utf8",
	"tls=skip-verify",
	"parseTime=true",
	fmt.Sprintf("loc=%s", url.QueryEscape("Asia/Shanghai")),
}

func ConnectMySQL() (*sql.DB, error) {
	config := utils.GetConfig()
	params := config.GetSubsSettings().Db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?%s",
		params.Username, params.Password, params.Url, params.Name, strings.Join(PARAMS, "&")))
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	// 如果超出最大连接数，等待可用的连接
	for db.Stats().OpenConnections >= config.GetSubsSettings().Db.MaxConn {
		time.Sleep(5 * time.Second)
	}
	return db, nil
}