package dao

import (
	"sync"
	"unsafe"
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

func (d *collectDao) AddSentCollect(txHash string, asset string, address string, amount float64) (int64, error) {
	params := make(map[string]interface {})
	params["address"] = address
	params["amount"] = amount
	params["asset"] = asset
	if txHash != "" {
		params["tx_hash"] = txHash
	}
	return insertPartsTemplate((*baseDao)(unsafe.Pointer(d)), "AddSentCollect", params)
}