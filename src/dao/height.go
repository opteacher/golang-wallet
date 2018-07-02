package dao

import (
	"sync"
	"database/sql"
	"utils"
	"unsafe"
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

func (d *heightDao) ChkOrAddAsset(asset string) (int64, error) {
	conds := []interface {} { asset }
	return saveTemplate((*baseDao)(unsafe.Pointer(d)),
		"GetHeight", "AddAsset", "",
		conds, conds, nil)
}

func (d *heightDao) GetHeight(asset string) (int64, error) {
	var result []map[string]interface {}
	var err error
	var conds = []interface {} { asset }
	bd :=(*baseDao)(unsafe.Pointer(d))
	if result, err = selectTemplate(bd, "GetHeight", conds); err != nil {
		return 0, err
	}
	if len(result) != 1 {
		return 0, utils.LogMsgEx(utils.WARNING, "获取高度失败", nil)
	}
	return result[0]["height"].(*sql.NullInt64).Int64, nil
}

func (d *heightDao) UpdateHeight(asset string, height uint64) (int64, error) {
	conds := []interface {} { asset }
	props := map[string]interface {} {
		"height": height,
	}
	return updatePartsTemplate((*baseDao)(unsafe.Pointer(d)), "UpdateHeight", conds, props)
}