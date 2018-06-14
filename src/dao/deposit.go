package dao

import (
	"sync"
	"database/sql"
	"databases"
	"log"
	"errors"
	"entities"
	"time"
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

func (dao *DepositDao) AddScannedDeposit(deposit entities.BaseDeposit) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var result sql.Result
	var insertSQL string
	var ok bool
	var useSQL = "AddScannedDeposit"
	if deposit.CreateTime.Year() > 1000 {
		useSQL = "AddDepositWithTime"
	}
	if insertSQL, ok = dao.sqls[useSQL]; !ok {
		err = errors.New("Cant find insert [deposit] table SQL")
		log.Println(err)
		return 0, err
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
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DepositDao) AddStableDeposit(deposit entities.BaseDeposit) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var insertSQL string
	var ok bool
	if deposit.CreateTime.Year() > 1000 {
		deposit.CreateTime = time.Now()
	}
	if insertSQL, ok = dao.sqls["AddStableDeposit"]; !ok {
		err = errors.New("Cant find insert [deposit] table SQL")
		log.Println(err)
		return 0, err
	}

	var params = []interface {} {
		deposit.TxHash,
		deposit.Address,
		deposit.Amount,
		deposit.Asset,
		deposit.Height,
		deposit.TxIndex,
		deposit.CreateTime,
		deposit.CreateTime,
	}
	var result sql.Result
	if result, err = db.Exec(insertSQL, params...); err != nil {
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *DepositDao) GetUnstableDeposit(asset string) ([]entities.TotalDeposit, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var selectSQL string
	var ok bool
	if selectSQL, ok = dao.sqls["GetUnstableDeposit"]; !ok {
		err = errors.New("Cant find select [deposit] table SQL")
		log.Println(err)
		return []entities.TotalDeposit {}, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, asset); err != nil {
		log.Println(err)
		return []entities.TotalDeposit {}, err
	}
	defer rows.Close()

	deposits := []entities.TotalDeposit {}
	for rows.Next() {
		var deposit entities.TotalDeposit
		var createTime = new(time.Time)
		var updateTime = new(time.Time)
		if err = rows.Scan([]interface {} {
			&deposit.Id,
			&deposit.TxHash,
			&deposit.Address,
			&deposit.Amount,
			&deposit.Asset,
			&deposit.Height,
			&deposit.TxIndex,
			&deposit.Status,
			createTime,
			updateTime,
		}...); err != nil {
			log.Println(err)
			continue
		}
		if createTime != nil {
			deposit.CreateTime = *createTime
		}
		if updateTime != nil {
			deposit.UpdateTime = *updateTime
		}
		deposits = append(deposits, deposit)
	}

	if err = rows.Err(); err != nil {
		return []entities.TotalDeposit {}, err
	}
	return deposits, nil
}

func (dao *DepositDao) DepositIntoStable(txHash string) (int64, error) {
	return 0, nil
}