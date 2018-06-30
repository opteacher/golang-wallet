package dao

import (
	"sync"
	"entities"
	"unsafe"
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

func (dao *withdrawDao) GetMaxId() (int, error) {
	return 0, nil
}

func (d *withdrawDao) NewWithdraw(withdraw entities.BaseWithdraw) (int64, error) {
	return insertTemplate((*baseDao)(unsafe.Pointer(d)), "NewWithdraw", []interface {} {
		withdraw.Address,
		withdraw.Amount,
		withdraw.Asset,
	})
}

func (d *withdrawDao) WithdrawIntoStable(txHash string) (int64, error) {
	return updateTemplate((*baseDao)(unsafe.Pointer(d)), "WithdrawIntoStable",
		[]interface {} { txHash }, nil)
}