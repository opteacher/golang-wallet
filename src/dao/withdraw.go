package dao

import (
	"sync"
	"entities"
	"database/sql"
	"databases"
	"utils"
)

type withdrawDao struct {
	baseDao
	sync.Once
}

var _withdrawDao *withdrawDao

func GetWithdrawDAO() *withdrawDao {
	if _withdrawDao == nil {
		_withdrawDao = new(withdrawDao)
		_withdrawDao.Once = sync.Once {}
		_withdrawDao.Once.Do(func() {
			_withdrawDao.create("withdraw")
		})
	}
	return _withdrawDao
}

func (dao *withdrawDao) GetAllUnstable() ([]entities.DatabaseWithdraw, error) {
	return []entities.DatabaseWithdraw {}, nil
}

func (dao *withdrawDao) NewWithdraw(withdraw entities.BaseWithdraw) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var insertSQL string
	var ok bool
	if insertSQL, ok = dao.sqls["NewWithdraw"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "NewWithdraw")
	}

	var result sql.Result
	if result, err = db.Exec(insertSQL, []interface {} {
		withdraw.Address,
		withdraw.Amount,
		withdraw.Asset,
	}...); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0012, err))
	}
	return result.RowsAffected()
}