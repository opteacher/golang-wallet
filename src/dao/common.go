package dao

import (
	"os"
	"log"
	"bufio"
	"strings"
	"io"
	"database/sql"
	"databases"
	"errors"
	"fmt"
)

type baseDao struct {
	sqls map[string]string
}

func (dao *baseDao) create(sqlFile string) error {
	var err error
	if dao.sqls, err = loadSQL(fmt.Sprintf("sql/%s.sql", sqlFile)); err != nil {
		log.Fatal(err)
	}

	var db *sql.DB
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var createSQL string
	var ok bool
	if createSQL, ok = dao.sqls["CreateTable"]; ok {
		if _, err = db.Exec(createSQL); err != nil {
			log.Fatal(err)
		}
	} else {
		err = errors.New(fmt.Sprintf("Cant find create [%s] table SQL", sqlFile))
		log.Println(err)
		return err
	}
	return nil
}

func loadSQL(sqlFile string) (map[string]string, error) {
	var file *os.File
	var err error
	if file, err = os.Open(sqlFile); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	sqlMap := map[string]string {}
	var key string
	var fmrKey string
	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if line == "" || line == "\n" {
			fmrKey = key
			key = ""
		} else if line[0] == '#' && len(key) != 0 {
			log.Println("SQL file structure has errors")
			break
		} else if line[0] == '#' {
			key = line[1:]
			key = strings.TrimSpace(key)
			sqlMap[key] = ""
		} else if len(key) == 0 {
			key = fmrKey
			sqlMap[key] += line
		} else {
			sqlMap[key] += line
		}

		if err == io.EOF {
			break
		}
	}
	return sqlMap, nil
}
