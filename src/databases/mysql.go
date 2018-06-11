package databases

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const DBUrl			= ""
const DBName		= "test"
const DBUserName	= "root"
const DBPassword	= "12345"

func Connect() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8&tls=skip-verify",
		DBUserName, DBPassword, DBUrl, DBName))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}