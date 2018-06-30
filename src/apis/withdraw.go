package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"entities"
	"strings"
	"dao"
	"utils"
	"services"
)

const revWithdraw = "^/api/withdraw/([A-Z]{3,})"

var wdRouteMap = map[string]api {
	revWithdraw:	{ doWithdraw, "POST" },
}

type withdrawReqBody struct {
	Id int			`json:"id"`
	Target string	`json:"target"`
	Value float64	`json:"value"`

}

func doWithdraw(w http.ResponseWriter, req *http.Request) []byte {
	var resp RespVO
	re := regexp.MustCompile(revWithdraw)
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
	var wdReqBody withdrawReqBody
	if err = json.Unmarshal(body, &wdReqBody); err != nil {
		utils.LogMsgEx(utils.WARNING, "转换请求体JSON对象错误：%v", err)
		resp.Code = 500
		resp.Msg = err.Error()
		ret, _ := json.Marshal(resp)
		return ret
	}

	// 参数判断
	var wdToSvc entities.BaseWithdraw
	if wdReqBody.Id == 0 {
		// 没有指定提币id，从数据库中挑选最大的id值
		if wdToSvc.Id, err = dao.GetWithdrawDAO().GetMaxId(); err != nil {
			utils.LogMsgEx(utils.WARNING, "从数据库获取提币ID错误：%v", err)
			resp.Code = 500
			resp.Msg = err.Error()
			ret, _ := json.Marshal(resp)
			return ret
		}
		wdToSvc.Id++
	} else {
		wdToSvc.Id = wdReqBody.Id
	}
	wdToSvc.Asset = strings.ToUpper(asset)
	if wdReqBody.Value == 0 {
		utils.LogMsgEx(utils.WARNING, "提币金额未指定", nil)
		resp.Code = 400
		resp.Msg = "提币金额未指定"
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		wdToSvc.Amount = wdReqBody.Value
	}
	if wdReqBody.Target == "" {
		utils.LogMsgEx(utils.WARNING, "提币目标地址不存在", nil)
		resp.Code = 400
		resp.Msg = "提币目标地址不存在"
		ret, _ := json.Marshal(resp)
		return ret
	} else {
		wdToSvc.Address = wdReqBody.Target
		wdToSvc.To = wdReqBody.Target
	}
	services.RevWithdrawSig <- wdToSvc
	return []byte(string(wdToSvc.Id))
}