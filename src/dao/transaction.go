package dao

import (
	"sync"
	"entities"
	"unsafe"
)

type transactionDao struct {
	baseDao
	sync.Once
}

var _transactionDao *transactionDao

func GetTransactionDAO() *transactionDao {
	if _transactionDao == nil {
		_transactionDao = new(transactionDao)
		_transactionDao.Once = sync.Once {}
		_transactionDao.Once.Do(func() {
			_transactionDao.create("transaction")
		})
	}
	return _transactionDao
}

func (d *transactionDao) AddTransaction(tx entities.Transaction, oprInf string) (int64, error) {
	entity := make(map[string]interface {})
	entity["opr_info"] = oprInf
	if tx.TxHash != "" {
		entity["tx_hash"] = tx.TxHash
	}
	if tx.BlockHash != "" {
		entity["block_hash"] = tx.BlockHash
	}
	if tx.From != "" {
		entity["`from`"] = tx.From
	}
	if tx.To != "" {
		entity["`to`"] = tx.To
	}
	if tx.Amount != 0 {
		entity["amount"] = tx.Amount
	}
	if tx.Asset != "" {
		entity["asset"] = tx.Asset
	}
	if tx.Height != 0 {
		entity["height"] = tx.Height
	}
	if tx.TxIndex != 0 {
		entity["tx_index"] = tx.TxIndex
	}
	if tx.CreateTime.Year() > 1990 {
		entity["create_time"] = tx.CreateTime
	}
	return insertPartsTemplate((*baseDao)(unsafe.Pointer(d)), "AddTransaction", entity)
}