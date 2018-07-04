package dao

import (
	"sync"
	"entities"
	"utils"
	"databases"
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

var _timeFormat = map[string]string {
	"ANSIC": time.ANSIC,
	"UnixDate": time.UnixDate,
	"RubyDate": time.RubyDate,
	"RFC822": time.RFC822,
	"RFC822Z": time.RFC822Z,
	"RFC850": time.RFC850,
	"RFC1123": time.RFC1123,
	"RFC1123Z": time.RFC1123Z,
	"RFC3339": time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen": time.Kitchen,
}

type processDao struct {
	baseDao
	sync.Once
}

var _processDao *processDao

func GetProcessDAO() *processDao {
	if _processDao == nil {
		_processDao = new(processDao)
		_processDao.Once = sync.Once {}
		_processDao.Once.Do(func() {})
	}
	return _processDao
}

func (d *processDao) SaveProcess(process *entities.DatabaseProcess) (int64, error) {
	var cli redis.Cmdable
	var err error
	if cli, err = databases.ConnectRedis(); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 42, err)
	}

	key := fmt.Sprintf("process_%s_%d", process.Asset, process.Id)

	if process.TxHash != "" {
		if err = cli.HSet(key, "tx_hash", process.TxHash).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置tx_hash失败：%v", err)
		}
	}
	if process.Asset != "" {
		if err = cli.HSet(key, "asset", process.Asset).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置asset失败：%v", err)
		}
	}
	if process.Type != "" {
		if err = cli.HSet(key, "type", process.Type).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置type失败：%v", err)
		}
	}
	if utils.StrArrayContains(entities.Processes, process.Process) {
		if err = cli.HSet(key, "process", process.Process).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置process失败：%v", err)
		}
	}
	if !process.Cancelable {
		if err = cli.HSet(key, "cancelable", 0).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置cancelable失败：%v", err)
		}
	}
	if process.Height != 0 {
		if err = cli.HSet(key, "height", process.Height).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置height失败：%v", err)
		}
	}
	if process.CurrentHeight != 0 {
		if err = cli.HSet(key, "current_height", process.CurrentHeight).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置current_height失败：%v", err)
		}
	}
	if process.CompleteHeight != 0 {
		if err = cli.HSet(key, "complete_height", process.CompleteHeight).Err(); err != nil {
			return 0, utils.LogMsgEx(utils.ERROR, "设置complete_height失败：%v", err)
		}
	}
	idTmFmt := utils.GetConfig().GetSubsSettings().Redis.TimeFormat
	var ok bool
	if idTmFmt, ok = _timeFormat[idTmFmt]; !ok {
		idTmFmt = _timeFormat["RFC3339"]
	}
	if err = cli.HSet(key, "last_update_time", time.Now().Format(idTmFmt)).Err(); err != nil {
		return 0, utils.LogMsgEx(utils.ERROR, "设置last_update_time失败：%v", err)
	}
	return 1, nil
}

func (d *processDao) QueryProcess(asset string, id int) (entities.DatabaseProcess, error) {
	var ret entities.DatabaseProcess
	var cli redis.Cmdable
	var err error
	if cli, err = databases.ConnectRedis(); err != nil {
		return ret, utils.LogIdxEx(utils.ERROR, 42, err)
	}

	key := fmt.Sprintf("process_%s_%d", asset, id)

	ret.Id					= id
	ret.TxHash, err			= cli.HGet(key, "tx_hash").Result()
	ret.Asset, err			= cli.HGet(key, "asset").Result()
	ret.Type, err			= cli.HGet(key, "type").Result()
	ret.Process, err		= cli.HGet(key, "process").Result()
	var strCancelable string
	strCancelable, err		= cli.HGet(key, "cancelable").Result()
	ret.Cancelable			= strCancelable == "true"
	ret.Height, err			= cli.HGet(key, "height").Uint64()
	ret.CurrentHeight, err	= cli.HGet(key, "current_height").Uint64()
	ret.CompleteHeight, err	= cli.HGet(key, "complete_height").Uint64()
	var strLstUpdTm string
	strLstUpdTm, err		= cli.HGet(key, "last_update_time").Result()
	idTmFmt := utils.GetConfig().GetSubsSettings().Redis.TimeFormat
	var ok bool
	if idTmFmt, ok = _timeFormat[idTmFmt]; !ok {
		idTmFmt = _timeFormat["RFC3339"]
	}
	ret.LastUpdateTime, err	= time.Parse(idTmFmt, strLstUpdTm)
	if err != nil {
		return ret, utils.LogMsgEx(utils.ERROR, "获取进度失败：%v", err)
	} else {
		return ret, nil
	}
}

func (d *processDao) UpdateHeight(asset string, curHeight uint64) (int64, error) {
	var cli redis.Cmdable
	var err error
	if cli, err = databases.ConnectRedis(); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 42, err)
	}

	var keys []string
	if keys, err = cli.Keys("process_*").Result(); err != nil {
		return 0, utils.LogMsgEx(utils.ERROR, "获取所有键失败：%v", err)
	}

	numKeys := len(keys)
	for _, key := range keys {
		if err = cli.HSet(key, "current_height", curHeight).Err(); err != nil {
			utils.LogMsgEx(utils.ERROR, "更新交易：%s失败：%v", key, err)
			numKeys--
		}
	}
	return int64(numKeys), nil
}