package apis

import (
	"utils"
	"encoding/json"
	"io/ioutil"
	"rpcs"
	"fmt"
	"net/http"
	"regexp"
)

const transferPath		= "^/api/transfer/([A-Z]{3,})$"
const minerPath			= "^/api/miner/([A-Z]{3,})$"

var txRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodPost, transferPath):	transfer,
	fmt.Sprintf("%s %s", http.MethodPut, minerPath):		doMining,
}

type transactionReq struct {
	From string		`json:"from"`
	To string		`json:"to"`
	Amount float64	`json:"amount"`
}

func transfer(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(transferPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	// 参数解析
	var body []byte
	var err error
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		utils.LogMsgEx(utils.WARNING, "解析请求体错误：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	defer req.Body.Close()

	utils.LogMsgEx(utils.INFO, "收到交易请求：%s", string(body))

	var txReq transactionReq
	if err = json.Unmarshal(body, &txReq); err != nil {
		utils.LogIdxEx(utils.WARNING, 38, err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	rpc := rpcs.GetRPC(params[0])
	var txHash string
	tradePwd := utils.GetConfig().GetCoinSettings().TradePassword
	if txHash, err = rpc.SendTransaction(txReq.From, txReq.To, txReq.Amount, tradePwd); err != nil {
		utils.LogMsgEx(utils.ERROR, "发送交易失败：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	resp.Code = 200
	resp.Data = txHash
	ret, _ := json.Marshal(resp)
	return []byte(ret)
}

func doMining(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(transferPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	var body []byte
	var err error
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		utils.LogMsgEx(utils.WARNING, "解析请求体错误：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	defer req.Body.Close()

	reqBody := make(map[string]interface {})
	if err = json.Unmarshal(body, reqBody); err != nil {
		utils.LogIdxEx(utils.WARNING, 38, err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	miningSpeed := 1
	if oMiningSpeed, ok := reqBody["speed"]; ok {
		miningSpeed = int(oMiningSpeed.(int64))
	}
	rpc := rpcs.GetRPC(params[0])
	if oEnable, ok := reqBody["enable"]; ok {
		if res, err := rpc.EnableMining(oEnable.(bool), miningSpeed); err != nil {
			utils.LogMsgEx(utils.WARNING, "调整挖矿状态失败：%v", err)
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		} else {
			resp.Code = 200
			resp.Data = res
			ret, _ := json.Marshal(resp)
			return []byte(ret)
		}
	}
	resp.Code = 200
	ret, _ := json.Marshal(resp)
	return []byte(ret)
}
