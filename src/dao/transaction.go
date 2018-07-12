package dao

import (
	"sync"
	"entities"
	"unsafe"
	"utils"
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
	bd := (*baseDao)(unsafe.Pointer(d))
	var result []map[string]interface {}
	var err error
	if result, err = selectTemplate(bd, "CheckExists", []interface {} { oprInf }); err != nil {
		return 0, utils.LogMsgEx(utils.ERROR, "检测交易存在与否失败：%v", err)
	}
	if len(result) != 1 {
		return 0, utils.LogMsgEx(utils.ERROR, "检测结果不等于1", nil)
	}
	var ok bool
	var tmp interface {}
	if tmp, ok = result[0]["num"]; !ok {
		return 0, utils.LogMsgEx(utils.ERROR, "检测存在数量的num分量不存在", nil)
	}
	if tmp.(int64) != 0 {
		utils.LogMsgEx(utils.WARNING, "交易：%s已被记录", oprInf)
		return 0, nil
	}

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
	return insertPartsTemplate(bd, "AddTransaction", entity)
}