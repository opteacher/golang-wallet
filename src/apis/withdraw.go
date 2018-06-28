package apis

import (
	"net/http"
	"regexp"
	"encoding/json"
	"io/ioutil"
	"utils"
	"entities"
	"services"
)

type withdrawReq struct {
	Id int			`json:"id"`
	Value float64	`json:"value"`
	Target string	`json:"target"`
}

const withdrawPath = "^/api/withdraw/([A-Z]{3,})"

var wdRouteMap = map[string]api {
	withdrawPath:	{ doWithdraw, "POST" },
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
	coinName := params[0]

	var body []byte
	var err error
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		utils.LogIdxEx(utils.WARNING, 38, err)
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

	var wd entities.BaseWithdraw
	wd.Id = wdReq.Id
	wd.Amount = wdReq.Value
	wd.Address = wdReq.Target
	wd.Asset = coinName
	services.RevWithdrawSig <- wd

	return []byte("")
}