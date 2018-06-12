package databases

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"utils"
)

func ConnectMySQL() (*sql.DB, error) {
	config := utils.GetConfig()
	params := config.GetSubsSettings().Db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8&tls=skip-verify",
		params.Username, params.Password, params.Url, params.Name))
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}