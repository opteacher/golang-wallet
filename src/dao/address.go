package dao

import (
	"sync"
	"databases"
	"log"
	"errors"
	"database/sql"
)

type AddressDao struct {
	baseDao
	sync.Once
}

var _addressDao *AddressDao

func GetAddressDAO() *AddressDao {
	if _addressDao == nil {
		_addressDao = new(AddressDao)
		_addressDao.Once = sync.Once {}
		_addressDao.Once.Do(func() {
			_addressDao.create("address")
		})
	}
	return _addressDao
}

func (dao *AddressDao) newAddress(asset string, address string, inuse bool) (int64, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Println(err)
		return 0, err
	}

	var result sql.Result
	var insertSQL string
	var ok bool
	var useSQL = "NewAddress"
	if inuse {
		useSQL = "NewAddressInuse"
	}
	if insertSQL, ok = dao.sqls[useSQL]; ok {
		if result, err = db.Exec(insertSQL, asset, address); err != nil {
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

func (dao *AddressDao) NewAddress(asset string, address string) (int64, error) {
	return dao.newAddress(asset, address, false)
}

func (dao *AddressDao) NewAddressInuse(asset string, address string) (int64, error) {
	return dao.newAddress(asset, address, true)
}

func (dao *AddressDao) FindInuseByAsset(asset string) ([]string, error) {
	var db *sql.DB
	var err error
	if db, err = databases.ConnectMySQL(); err != nil {
		log.Println(err)
		return nil, err
	}

	var selectSQL string
	var ok bool
	if selectSQL, ok = dao.sqls["FindByAsset"]; !ok {
		err = errors.New("Cant find insert [address] table SQL")
		log.Println(err)
		return []string {}, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(selectSQL, asset); err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err = rows.Scan(&address); err != nil {
			log.Println(err)
			continue
		}
		addresses = append(addresses, address)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return addresses, nil
}