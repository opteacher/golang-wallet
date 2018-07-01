package databases

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"utils"
	"strings"
	"net/url"
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
	return db, nil
}