package dao

import (
	"os"
	"bufio"
	"strings"
	"io"
	"database/sql"
	"databases"
	"errors"
	"fmt"
	"utils"
)

type baseDao struct {
	sqls map[string]string
}

func (dao *baseDao) create(sqlFile string) error {
	var err error
	if dao.sqls, err = loadSQL(fmt.Sprintf("sql/%s.sql", sqlFile)); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0015, err))
	}

	var db *sql.DB
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var createSQL string
	var ok bool
	if createSQL, ok = dao.sqls["CreateTable"]; ok {
		if _, err = db.Exec(createSQL); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 0016, err))
		}
	} else {
		return utils.LogIdxEx(utils.ERROR, 0011, errors.New("CreateTable"))
	}
	return nil
}

func loadSQL(sqlFile string) (map[string]string, error) {
	var file *os.File
	var err error
	if file, err = os.Open(sqlFile); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0017, err))
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
			panic(utils.LogIdxEx(utils.ERROR, 18, err))
		}

		if line == "" || line == "\n" {
			fmrKey = key
			key = ""
		} else if line[0] == '#' && len(key) != 0 {
			utils.LogIdxEx(utils.ERROR, 19, errors.New("SQL标题要以#开头"))
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
