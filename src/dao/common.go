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
	"reflect"
)

type baseDao struct {
	sqls map[string]string
}

func (dao *baseDao) create(sqlFile string) error {
	var err error
	if dao.sqls, err = loadSQL(fmt.Sprintf("sql/%s.sql", sqlFile)); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 15, err))
	}

	var db *sql.DB
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}

	var createSQL string
	var ok bool
	if createSQL, ok = dao.sqls["CreateTable"]; ok {
		if _, err = db.Exec(createSQL); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 16, err))
		}
	} else {
		return utils.LogIdxEx(utils.ERROR, 11, errors.New("CreateTable"))
	}
	return nil
}

func loadSQL(sqlFile string) (map[string]string, error) {
	var file *os.File
	var err error
	if file, err = os.Open(sqlFile); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 17, err))
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

func insertTemplate(d *baseDao, sqlName string, props []interface {}) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var insertSQL string
	var ok bool
	if insertSQL, ok = d.sqls[sqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var result sql.Result
	if result, err = db.Exec(insertSQL, props...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 12, err))
	}
	return result.RowsAffected()
}

func saveTemplate(d *baseDao, selSqlName string, istSqlName string, updSqlName string, conds []interface{}, props []interface{}, keys []string) (int64, error) {

	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var selectSQL string
	var ok bool
	if selectSQL, ok = d.sqls[selSqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, selSqlName)
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, conds...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 13, err))
	}
	defer rows.Close()

	var result sql.Result
	if !rows.Next() {
		if len(keys) != 0 {
			if len(props) != len(keys) {
				panic(utils.LogMsgEx(utils.ERROR, "键值无法一一对应", nil))
			}

			params := make(map[string]interface {})
			for i, k := range keys {
				params[k] = props[i]
			}
			return insertPartsTemplate(d, istSqlName, params)
		} else {
			var insertSQL string
			var ok bool
			if insertSQL, ok = d.sqls[istSqlName]; !ok {
				return 0, utils.LogIdxEx(utils.ERROR, 11, istSqlName)
			}

			if result, err = db.Exec(insertSQL, props...); err != nil {
				panic(utils.LogIdxEx(utils.ERROR, 12, err))
			}
			return result.RowsAffected()
		}
	} else {
		if updSqlName == "" {
			return 0, nil
		}
		if len(props) != len(keys) {
			panic(utils.LogMsgEx(utils.ERROR, "键值无法一一对应", nil))
		}

		var updateSQL string
		if updateSQL, ok = d.sqls[updSqlName]; !ok {
			return 0, utils.LogIdxEx(utils.ERROR, 11, updSqlName)
		}

		var content []string
		for _, key := range keys {
			content = append(content, key + "=?")
		}
		updateSQL = fmt.Sprintf(updateSQL, strings.Join(content, ","))

		if result, err = db.Exec(updateSQL, append(props, conds...)...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 12, err))
		}
	}
	return result.RowsAffected()
}

func selectTemplate(d *baseDao, sqlName string, conds []interface {}) ([]map[string]interface {}, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var selectSQL string
	var ok bool
	if selectSQL, ok = d.sqls[sqlName]; !ok {
		return nil, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, conds...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 13, err))
	}
	defer rows.Close()

	var result []map[string]interface {}
	for rows.Next() {
		var entity = make(map[string]interface {})
		var colTyps []*sql.ColumnType
		if colTyps, err = rows.ColumnTypes(); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 14, err))
		}
		var params []interface {}
		for _, colTyp := range colTyps {
			entity[colTyp.Name()] = reflect.New(colTyp.ScanType()).Interface()
			params = append(params, entity[colTyp.Name()])
		}
		if err = rows.Scan(params...); err != nil {
			utils.LogIdxEx(utils.ERROR, 14, err)
			continue
		}
		result = append(result, entity)
	}

	if err = rows.Err(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 14, err))
	}
	return result, nil
}

func selectPartsTemplate(d *baseDao, sqlName string, conds map[string]interface{}) ([]map[string]interface {}, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var selectSQL string
	var ok bool
	if selectSQL, ok = d.sqls[sqlName]; !ok {
		return nil, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var rows *sql.Rows
	var keys []string
	var vals []interface {}
	for key, val := range conds {
		keys = append(keys, fmt.Sprintf("%s=?", key))
		vals = append(vals, val)
	}
	selectSQL = fmt.Sprintf(selectSQL, strings.Join(keys, " AND "))
	if rows, err = db.Query(selectSQL, vals...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 13, err))
	}
	defer rows.Close()

	var result []map[string]interface {}
	for rows.Next() {
		var entity = make(map[string]interface {})
		var colTyps []*sql.ColumnType
		if colTyps, err = rows.ColumnTypes(); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 14, err))
		}
		var params []interface {}
		for _, colTyp := range colTyps {
			entity[colTyp.Name()] = reflect.New(colTyp.ScanType()).Interface()
			params = append(params, entity[colTyp.Name()])
		}
		if err = rows.Scan(params...); err != nil {
			utils.LogIdxEx(utils.ERROR, 14, err)
			continue
		}
		result = append(result, entity)
	}

	if err = rows.Err(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 14, err))
	}
	return result, nil
}

func updateTemplate(d *baseDao, sqlName string, conds []interface {}, props []interface {}) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var updateSQL string
	var ok bool
	if updateSQL, ok = d.sqls[sqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var result sql.Result
	if props == nil || len(props) == 0 {
		if result, err = db.Exec(updateSQL, conds...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 21, err))
		}
		return result.RowsAffected()
	} else {
		if result, err = db.Exec(updateSQL, append(props, conds...)...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 21, err))
		}
		return result.RowsAffected()
	}
}

func updatePartsTemplate(d *baseDao, sqlName string, conds []interface {}, props map[string]interface{}) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var updateSQL string
	var ok bool
	if updateSQL, ok = d.sqls[sqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var result sql.Result
	if props == nil || len(props) == 0 {
		if result, err = db.Exec(updateSQL, conds...); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 21, err))
		}
		return result.RowsAffected()
	}

	var content []string
	var values []interface {}
	for k, v := range props {
		content = append(content, k + "=?")
		values = append(values, v)
	}

	updateSQL = fmt.Sprintf(updateSQL, strings.Join(content, ","))
	if result, err = db.Exec(updateSQL, append(values, conds...)...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 21, err))
	}
	return result.RowsAffected()
}

func insertPartsTemplate(d *baseDao, sqlName string, props map[string]interface {}) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var insertSQL string
	var ok bool
	if insertSQL, ok = d.sqls[sqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	var keys []string
	var vals []string
	var propts []interface {}
	for k, v := range props {
		keys = append(keys, k)
		vals = append(vals, "?")
		propts = append(propts, v)
	}

	var result sql.Result
	insertSQL = fmt.Sprintf(insertSQL, strings.Join(keys, ","), strings.Join(vals, ","))
	if result, err = db.Exec(insertSQL, propts...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 12, err))
	}
	return result.RowsAffected()
}

func deleteTemplate(d *baseDao, sqlName string, conds []interface {}) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 10, err))
	}
	defer db.Close()

	var deleteSQL string
	var ok bool
	if deleteSQL, ok = d.sqls[sqlName]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 11, sqlName)
	}

	if result, err := db.Exec(deleteSQL, conds...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 21, err))
	} else {
		return result.RowsAffected()
	}
}