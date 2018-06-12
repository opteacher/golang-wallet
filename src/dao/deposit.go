package dao

import (
	"sync"
	"database/sql"
	"databases"
	"log"
	"errors"
	"entities"
)

type DepositDao struct {
	baseDao
	sync.Once
}

var _depositDao *DepositDao

func GetDepositDAO() *DepositDao {
	if _depositDao == nil {
		_depositDao = new(DepositDao)
		_depositDao.Once = sync.Once {}
		_depositDao.Once.Do(func() {
			_depositDao.create("deposit")
		})
	}
	return _depositDao
}

func (dao *DepositDao) AddScannedDeposit(deposit entities.Deposit) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Println(err)
		return 0, err
	}

	var result sql.Result
	var insertSQL string
	var ok bool
	if insertSQL, ok = dao.sqls["AddScannedDeposit"]; ok {
		if result, err = db.Exec(insertSQL,
			deposit.TxHash,
			deposit.Address,
			deposit.Amount,
			deposit.Asset,
			deposit.Height,
			deposit.TxIndex,
		); err != nil {
			log.Println(err)
			return 0, err
		}
	} else {
		err = errors.New("Cant find insert [address] table SQL")
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}