package dao

import (
	"sync"
	"database/sql"
	"databases"
	"utils"
	"fmt"
)

type collectDao struct {
	baseDao
	sync.Once
}

var _collectDao *collectDao

func GetCollectDAO() *collectDao {
	if _collectDao == nil {
		_collectDao = new(collectDao)
		_collectDao.Once = sync.Once {}
		_collectDao.Once.Do(func() {
			_collectDao.create("collect")
		})
	}
	return _collectDao
}

func (dao *collectDao) AddSentCollect(txHash string, asset string, address string, amount float64) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var result sql.Result
	var insertSQL string
	var ok bool
	if insertSQL, ok = dao.sqls["AddSentCollect"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "AddSentCollect")
	}

	txHashCol := "tx_hash,"
	txHashVal := "?,"
	params := []interface {} {
		txHash, address, amount, asset,
	}
	if txHash == "" {
		txHashCol = ""
		txHashVal = ""
		params = params[1:]
	}
	insertSQL = fmt.Sprintf(insertSQL, txHashCol, txHashVal)

	if result, err = db.Exec(insertSQL, params...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0012, err))
	}
	return result.RowsAffected()
}