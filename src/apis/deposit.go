package apis

import (
	"net/http"
	"regexp"
	"rpcs"
	"encoding/json"
	"dao"
	"entities"
	"fmt"
)

const newAddressPath	= "^/api/deposit/([A-Z]{3,})/address$"
const getHeightPath		= "^/api/deposit/([A-Z]{3,})/height$"
const getDepositsPath	= "^/api/deposit/([A-Z]{3,})$"

var dpRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodGet, newAddressPath):	newAddress,
	fmt.Sprintf("%s %s", http.MethodGet, getHeightPath):	queryHeight,
	fmt.Sprintf("%s %s", http.MethodGet, getDepositsPath):	getDeposits,
}

func newAddress(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(newAddressPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	asset := params[0]
	rpc := rpcs.GetRPC(asset)

	var address string
	var err error
	if address, err = rpc.GetNewAddress(); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	// 保存进充值地址数据库
	if _, err = dao.GetAddressDAO().NewAddressInuse(asset, address); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = address
	ret, _ := json.Marshal(resp)
	return ret
}
func queryHeight(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(getHeightPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	coinName := params[0]
	rpc := rpcs.GetRPC(coinName)

	var height uint64
	var err error
	if height, err = rpc.GetCurrentHeight(); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = height
	ret, _ := json.Marshal(resp)
	return ret
}
func getDeposits(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(getDepositsPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	conds := make(map[string]interface {})
	conds["asset"] = params[0]
	var result []entities.DatabaseDeposit
	var err error
	if txHash := req.Form.Get("tx_hash"); txHash != "" {
		conds["tx_hash"] = txHash

		// txhash是唯一的，所以指定的话，直接返回
		if result ,err = dao.GetDepositDAO().GetDeposits(conds); err != nil {
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		}
		resp.Code = 200
		resp.Data = result
		ret, _ := json.Marshal(resp)
		return ret
	}
	if address := req.Form.Get("address"); address != "" {
		conds["address"] = address
	}

	if result ,err = dao.GetDepositDAO().GetDeposits(conds); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	resp.Code = 200
	resp.Data = result
	ret, _ := json.Marshal(resp)
	return ret
}