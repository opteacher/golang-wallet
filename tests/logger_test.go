package tests

import (
	"utils"
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	utils.LogMsgEx(utils.ERROR, "测试1：abcd", nil)
	a := 2000
	utils.LogMsgEx(utils.WARNING, "测试2：%d", a)

	utils.LogIdxEx(utils.INFO, 0, nil)

	err := errors.New("error test")
	utils.LogIdxEx(utils.DEBUG, 0, err)
}