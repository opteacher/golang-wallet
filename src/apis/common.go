package apis

import (
	"net/http"
	"fmt"
	"utils"
)

const WithdrawPath = "/api/withdraw"
const WpLen = len(WithdrawPath)
const DepositPath = "/api/deposit"
const DpLen = len(DepositPath)

func RootHandler(w http.ResponseWriter, req *http.Request) {
	utils.LogMsgEx(utils.INFO, "%s\t\t%s", req.Method, req.RequestURI)
	switch {
	case len(req.RequestURI) >= WpLen && req.RequestURI[:WpLen] == WithdrawPath:
		WithdrawHandler(w, req)
	case len(req.RequestURI) >= DpLen && req.RequestURI[:DpLen] == DepositPath:
		fmt.Println(req.RequestURI[DpLen:])
	}
}