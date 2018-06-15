package dao

import (
	"sync"
	"database/sql"
	"databases"
	"utils"
)

type heightDao struct {
	baseDao
	sync.Once
}

var _heightDao *heightDao

func GetHeightDAO() *heightDao {
	if _heightDao == nil {
		_heightDao = new(heightDao)
		_heightDao.Once = sync.Once {}
		_heightDao.Once.Do(func() {
			_heightDao.create("height")
		})
	}
	return _heightDao
}

func (dao *heightDao) ChkOrAddAsset(asset string) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var checkSQL string
	var ok bool
	if checkSQL, ok = dao.sqls["ChkExsAsset"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "ChkExsAsset")
	}

	var rows *sql.Rows
	if rows, err = db.Query(checkSQL, asset); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0013, err))
	}
	defer rows.Close()

	var num int
	if rows.Next() {
		if err = rows.Scan(&num); err != nil {
			return 0, utils.LogIdxEx(utils.ERROR, 0014, err)
		}
	}

	if num > 0 {
		return 0, nil
	}

	var insertSQL string
	if insertSQL, ok = dao.sqls["AddAsset"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "AddAsset")
	}

	var result sql.Result
	if result, err = db.Exec(insertSQL, asset); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0012, err))
	}
	return result.RowsAffected()
}

func (dao *heightDao) GetHeight(asset string) (uint64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var selectSQL string
	var ok bool
	if selectSQL, ok = dao.sqls["GetHeight"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "GetHeight")
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, asset); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0013, err))
	}
	defer rows.Close()

	var height uint64
	if rows.Next() {
		if err = rows.Scan(&height); err != nil {
			return 0, utils.LogIdxEx(utils.ERROR, 0014, err)
		} else {
			return height, nil
		}
	}

	if err = rows.Err(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0014, err))
	}
	return 0, utils.LogMsgEx(utils.WARNING, "无法找到指定币种的高度：%s", asset)
}

func (dao *heightDao) UpdateHeight(asset string, height uint64) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0010, err))
	}

	var updateSQL string
	var ok bool
	if updateSQL, ok = dao.sqls["UpdateHeight"]; !ok {
		return 0, utils.LogIdxEx(utils.ERROR, 0011, "UpdateHeight")
	}

	var result sql.Result
	if result, err = db.Exec(updateSQL, height, asset); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0021, err))
	}
	return result.RowsAffected()
}