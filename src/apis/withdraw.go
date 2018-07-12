package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"utils"
	"entities"
	"services"
	"dao"
	"strings"
	"fmt"
)

type withdrawReq struct {
	Id int			`json:"id"`
	Value float64	`json:"value"`
	Target string	`json:"target"`
}

const withdrawPath = "^/api/withdraw/([A-Z]{3,})"

var wdRouteMap = map[string]interface {} {
	fmt.Sprintf("%s %s", http.MethodPost, withdrawPath):	doWithdraw,
	fmt.Sprintf("%s %s", http.MethodGet, withdrawPath):		getWithdraw,
}

func doWithdraw(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(withdrawPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}
	asset := params[0]

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

	utils.LogMsgEx(utils.INFO, "收到提币请求：%s", string(body))

	var wdReq withdrawReq
	if err = json.Unmarshal(body, &wdReq); err != nil {
		utils.LogIdxEx(utils.WARNING, 38, err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	// 参数判断
	var wdToSvc entities.BaseWithdraw
	if wdReq.Id == 0 {
		// 没有指定提币id，从数据库中挑选最大的id值
		asset := utils.GetConfig().GetCoinSettings().Name
		if wdToSvc.Id, err = dao.GetWithdrawDAO().GetAvailableId(asset); err != nil {
			utils.LogMsgEx(utils.WARNING, "从数据库获取提币ID错误：%v", err)
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		}
	}

	var exist bool
	if exist, err = dao.GetWithdrawDAO().CheckExistsById(wdReq.Id); err != nil {
		utils.LogMsgEx(utils.WARNING, "从数据库检查提币ID错误：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}
	if exist {
		errStr := fmt.Sprintf("收到重复的提币请求，Id：%d", wdReq.Id)
		utils.LogMsgEx(utils.WARNING, errStr, nil)
		resp.Code = 500
		resp.Msg = errStr
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		wdToSvc.Id = wdReq.Id
	}
	wdToSvc.Asset = strings.ToUpper(asset)
	if wdReq.Value == 0 {
		utils.LogMsgEx(utils.WARNING, "提币金额未指定", nil)
		resp.Code = 400
		resp.Msg = "提币金额未指定"
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		wdToSvc.Amount = wdReq.Value
	}
	if wdReq.Target == "" {
		utils.LogMsgEx(utils.WARNING, "提币目标地址不存在", nil)
		resp.Code = 400
		resp.Msg = "提币目标地址不存在"
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		wdToSvc.Address = wdReq.Target
		wdToSvc.To = wdReq.Target
	}
	services.RevWithdrawSig <- wdToSvc

	resp.Code = 200
	resp.Data = wdToSvc.Id
	ret, _ := json.Marshal(resp)
	return []byte(ret)
}

func getWithdraw(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(withdrawPath)
	params := re.FindStringSubmatch(req.RequestURI)[1:]
	if len(params) == 0 {
		resp.Code = 500
		resp.Msg = "需要指定币种的名字"
		ret, _ := json.Marshal(resp)
		return ret
	}

	conds := make(map[string]interface {})
	conds["asset"] = params[0]
	var result []entities.DatabaseWithdraw
	var err error
	if txHash := req.Form.Get("tx_hash"); txHash != "" {
		conds["tx_hash"] = txHash

		// txhash是唯一的，所以指定的话，直接返回
		if result ,err = dao.GetWithdrawDAO().GetWithdraws(conds); err != nil {
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

	resp.Code = 200
	resp.Data = result
	ret, _ := json.Marshal(resp)
	return ret
}