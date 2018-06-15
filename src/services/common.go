package services

import (
	"dao"
	"entities"
	"utils"
)

const (
	DESTORY = iota
	STOP
	NONE
	CREATE
	INIT
	START
)

func TxIntoStable(deposit *entities.BaseDeposit, insert bool) error {
	var err error
	if insert {
		if _, err = dao.GetDepositDAO().AddStableDeposit(deposit); err != nil {
			return utils.LogMsgEx(utils.ERROR, "添加充币记录失败：%v", err)
		}
	} else {
		if _, err = dao.GetDepositDAO().DepositIntoStable(deposit.TxHash); err != nil {
			return utils.LogMsgEx(utils.ERROR, "更新充币记录失败：%v", err)
		}
	}

	var process entities.DatabaseProcess
	process.TxHash = deposit.TxHash
	process.Process = entities.FINISH
	if _, err = dao.GetProcessDAO().SaveProcess(&process); err != nil {
		return utils.LogMsgEx(utils.ERROR, "插入/更新进度表失败：%v", err)
	}

	utils.LogMsgEx(utils.INFO, "交易充值完成：%v", deposit)
	return nil
}

var toNotifySig = make(chan entities.BaseDeposit)