package dao

import (
	"sync"
	"database/sql"
	"unsafe"
)

type addressDao struct {
	baseDao
	sync.Once
}

var _addressDao *addressDao

func GetAddressDAO() *addressDao {
	if _addressDao == nil {
		_addressDao = new(addressDao)
		_addressDao.Once = sync.Once {}
		_addressDao.Once.Do(func() {
			_addressDao.create("address")
		})
	}
	return _addressDao
}

func (d *addressDao) newAddress(asset string, address string, inuse bool) (int64, error) {
	useSQL := "NewAddress"
	if inuse { useSQL = "NewAddressInuse" }
	props := []interface {} { asset, address }
	return insertTemplate((*baseDao)(unsafe.Pointer(d)), useSQL, props)
}

func (d *addressDao) NewAddress(asset string, address string) (int64, error) {
	return d.newAddress(asset, address, false)
}

func (d *addressDao) NewAddressInuse(asset string, address string) (int64, error) {
	return d.newAddress(asset, address, true)
}

func (d *addressDao) FindInuseByAsset(asset string) ([]string, error) {
	conds := []interface {} { asset }
	var result []map[string]interface {}
	var err error
	if result, err = selectTemplate((*baseDao)(unsafe.Pointer(d)), "FindByAsset", conds); err != nil {
		return []string {}, err
	}
	var ret []string
	for _, entity := range result {
		ret = append(ret, string(*entity["address"].(*sql.RawBytes)))
	}
	return ret, nil
}