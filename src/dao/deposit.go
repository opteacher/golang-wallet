package dao

import (
	"sync"
	"database/sql"
	"databases"
	"entities"
	"time"
	"utils"
)

type depositDao struct {
	baseDao
	sync.Once
}

var _depositDao *depositDao

func GetDepositDAO() *depositDao {
	if _depositDao == nil {
		_depositDao = new(depositDao)
		_depositDao.Once = sync.Once {}
		_depositDao.Once.Do(func() {
			_depositDao.create("deposit")
		})
	}
	return _depositDao
}

func (dao *depositDao) AddScannedDeposit(deposit *entities.BaseDeposit) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var result sql.Result
	var insertSQL string
	var ok bool
	var useSQL = "AddScannedDeposit"
	if deposit.CreateTime.Year() > 1000 {
		useSQL = "AddDepositWithTime"
	}
	if insertSQL, ok = dao.sqls[useSQL]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, useSQL)
	}

	var params = []interface {} {
		deposit.TxHash,
		deposit.Address,
		deposit.Amount,
		deposit.Asset,
		deposit.Height,
		deposit.TxIndex,
	}
	if useSQL == "AddDepositWithTime" {
		params = append(params, deposit.CreateTime)
	}
	if result, err = db.Exec(insertSQL, params...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0012, err))
	}
	return result.RowsAffected()
}

func (dao *depositDao) AddStableDeposit(deposit *entities.BaseDeposit) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var insertSQL string
	var ok bool
	if deposit.CreateTime.Year() < 1000 {
		deposit.CreateTime = time.Now()
	}
	if insertSQL, ok = dao.sqls["AddStableDeposit"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "AddStableDeposit")
	}

	var params = []interface {} {
		deposit.TxHash,
		deposit.Address,
		deposit.Amount,
		deposit.Asset,
		deposit.Height,
		deposit.TxIndex,
		deposit.CreateTime,
	}
	var result sql.Result
	if result, err = db.Exec(insertSQL, params...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0012, err))
	}
	return result.RowsAffected()
}

func (dao *depositDao) GetUnstableDeposit(asset string) ([]entities.BaseDeposit, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var selectSQL string
	var ok bool
	if selectSQL, ok = dao.sqls["GetUnstableDeposit"]; !ok {
		return nil, utils.LogIdxEx(utils.ERROR, 0011, "GetUnstableDeposit")
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, asset); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0013, err))
	}
	defer rows.Close()

	deposits := []entities.BaseDeposit {}
	for rows.Next() {
		var deposit entities.BaseDeposit
		if err = rows.Scan([]interface {} {
			&deposit.TxHash,
			&deposit.Address,
			&deposit.Amount,
			&deposit.Asset,
			&deposit.Height,
			&deposit.TxIndex,
		}...); err != nil {
			utils.LogIdxEx(utils.ERROR, 0014, err)
			continue
		}
		deposits = append(deposits, deposit)
	}

	if err = rows.Err(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0014, err))
	}
	return deposits, nil
}

func (dao *depositDao) DepositIntoStable(txHash string) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var updateSQL string
	var ok bool
	if updateSQL, ok = dao.sqls["DepositIntoStable"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "DepositIntoStable")
	}

	var result sql.Result
	if result, err = db.Exec(updateSQL, txHash); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0021, err))
	}
	return result.RowsAffected()
}