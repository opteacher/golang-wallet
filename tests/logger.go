package main

import (
	"utils"
	"errors"
)

func main() {
	utils.LogMsgEx(utils.ERROR, "测试1：abcd", nil)
	a := 2000
	utils.LogMsgEx(utils.WARNING, "测试2：%d", a)

	utils.LogIdxEx(utils.INFO, 0, nil)

	utils.EnableDebugLog(true)

	err := errors.New("error test")
	utils.LogIdxEx(utils.DEBUG, 0, err)
}