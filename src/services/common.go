package services

import (
	"dao"
	"entities"
	"utils"
	"unsafe"
)

const (
	DESTORY = iota
	STOP
	NONE
	CREATE
	INIT
	START
)

var ServiceStatus = []string {
	"DESTORY", "STOP", "NONE", "CREATE", "INIT", "START",
}

type BaseService struct {
	status utils.Status
	name string
}

func (service *BaseService) create() error {
	service.status.Init([]int { DESTORY, CREATE, INIT, START, STOP })
	return nil
}

func (service *BaseService) Init() {
	service.status.TurnTo(INIT)
}

func (service *BaseService) Start() {
	service.status.TurnTo(START)
}

func (service *BaseService) Stop() {
	if service.status.Current() == START {
		service.status.TurnTo(STOP)
	} else {
		service.status.TurnTo(DESTORY)
	}
}

func (service *BaseService) IsInit() bool {
	return service.status.Current() >= INIT
}

func (service *BaseService) IsDestroy() bool {
	return service.status.Current() == DESTORY
}

func (service *BaseService) CurrentStatus() int {
	return service.status.Current()
}

func (service *BaseService) Name() string {
	return service.name
}

func GetInitedServices() []*BaseService {
	var svcs []*BaseService
	if GetWithdrawService().IsInit() {
		svcs = append(svcs, (*BaseService)(unsafe.Pointer(GetWithdrawService())))
	}
	if GetCollectService().IsInit() {
		svcs = append(svcs, (*BaseService)(unsafe.Pointer(GetCollectService())))
	}
	if GetNotifyService().IsInit() {
		svcs = append(svcs, (*BaseService)(unsafe.Pointer(GetNotifyService())))
	}
	if GetDepositService().IsInit() {
		svcs = append(svcs, (*BaseService)(unsafe.Pointer(GetDepositService())))
	}
	return svcs
}

func TxIntoStable(txHash string, curHeight uint64) error {
	var err error
	if _, err = dao.GetDepositDAO().DepositIntoStable(txHash); err != nil {
		return utils.LogMsgEx(utils.ERROR, "更新充币记录失败：%v", err)
	}
	if _, err = dao.GetWithdrawDAO().WithdrawIntoStable(txHash); err != nil {
		return utils.LogMsgEx(utils.ERROR, "更新充币记录失败：%v", err)
	}

	var process entities.DatabaseProcess
	process.TxHash = txHash
	process.Height = curHeight
	process.Process = entities.FINISH
	if _, err = dao.GetProcessDAO().SaveProcess(&process); err != nil {
		return utils.LogMsgEx(utils.ERROR, "插入/更新进度表失败：%v", err)
	}

	utils.LogMsgEx(utils.INFO, "交易完成：%s", txHash)
	return nil
}

var toNotifySig = make(chan entities.Transaction)
var revWithdrawSig = make(chan entities.BaseWithdraw)