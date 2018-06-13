package dao

import (
	"sync"
	"database/sql"
	"databases"
	"log"
	"errors"
)

type HeightDao struct {
	baseDao
	sync.Once
}

var _heightDao *HeightDao

func GetHeightDAO() *HeightDao {
	if _heightDao == nil {
		_heightDao = new(HeightDao)
		_heightDao.Once = sync.Once {}
		_heightDao.Once.Do(func() {
			_heightDao.create("height")
		})
	}
	return _heightDao
}

func (dao *HeightDao) ChkOrAddAsset(asset string) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var checkSQL string
	var ok bool
	if checkSQL, ok = dao.sqls["ChkExsAsset"]; !ok {
		err = errors.New("Cant find insert [height] table SQL")
		log.Println(err)
		return 0, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(checkSQL, asset); err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()

	var num int
	if rows.Next() {
		if err = rows.Scan(&num); err != nil {
			log.Println(err)
			return 0, err
		}
	}

	if num > 0 {
		return 0, nil
	}

	var insertSQL string
	if insertSQL, ok = dao.sqls["AddAsset"]; !ok {
		err = errors.New("Cant find insert [height] table SQL")
		log.Println(err)
		return 0, err
	}

	var result sql.Result
	if result, err = db.Exec(insertSQL, asset); err != nil {
		log.Println(err)
		return 0, err
	}
	return result.RowsAffected()
}

func (dao *HeightDao) GetHeight(asset string) (uint64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Fatal(err)
	}

	var selectSQL string
	var ok bool
	if selectSQL, ok = dao.sqls["GetHeight"]; !ok {
		err = errors.New("Cant find insert [height] table SQL")
		log.Println(err)
		return 0, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, asset); err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()

	var height uint64
	if rows.Next() {
		if err = rows.Scan(&height); err != nil {
			return 0, err
		} else {
			return height, nil
		}
	}

	if err = rows.Err(); err != nil {
		return 0, err
	}
	err = errors.New("Cant find identified asset")
	return 0, nil
}