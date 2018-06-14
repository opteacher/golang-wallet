package rpcs

import (
	"sync"
	"utils"
	"entities"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	"fmt"
	"errors"
	"strconv"
	"math/rand"
	"time"
	"math/big"
	"math"
)

type eth struct {
	sync.Once
	coinName string
	callUrl string
	decimal int
	Stable int
}

var __eth *eth

func GetEth() *eth {
	if __eth == nil {
		__eth = new(eth)
		__eth.Once = sync.Once {}
		__eth.Once.Do(func() {
			__eth.create()
		})
	}
	return __eth
}

func (rpc *eth) create() {
	setting := utils.GetConfig().GetCoinSettings()
	rpc.coinName 	= setting.Name
	rpc.callUrl		= setting.Url
	rpc.decimal		= setting.Decimal
	rpc.Stable		= setting.Stable
}

type EthSucceedResp struct {
	JsonRpc string		`json:jsonrpc`
	Id string			`json:id`
	Result interface{}	`json:result`
}

type EthFailedResp struct {
	JsonRpc string		`json:jsonrpc`
	Id string			`json:id`
	Error struct {
		Code int		`json:code`
		Message string	`json:message`
	}					`json:error`
}

func (rpc *eth) sendRequest(method string, params []interface {}, id string) (EthSucceedResp, error)  {
	reqBody := RequestBody { method, params, id }
	reqStr, err := json.Marshal(reqBody)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0022, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Request body: %s", reqStr), nil)

	reqBuf := bytes.NewBuffer([]byte(reqStr))
	res, err := http.Post(rpc.callUrl, "application/json", reqBuf)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0024, err))
	}
	defer res.Body.Close()

	bodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 0025, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Response body: %s", bodyStr), nil)

	var resBody EthSucceedResp
	if err = json.Unmarshal(bodyStr, &resBody); err != nil {
		var resError EthFailedResp
		if err = json.Unmarshal(bodyStr, &resError); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 0023, err))
		} else {
			return resBody, utils.LogIdxEx(utils.ERROR, 0026, errors.New(resError.Error.Message))
		}
	}
	return resBody, nil
}

func (rpc *eth) GetTransactions(height uint, addresses []string) ([]entities.BaseDeposit, error) {
	var err error
	// 发送请求获取指定高度的块
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%d", rand.Intn(1000))
	params := []interface{} { "0x" + strconv.FormatUint(uint64(height), 16), true }
	if resp, err = rpc.sendRequest("eth_getBlockByNumber", params, id); err != nil {
		return nil, utils.LogIdxEx(utils.ERROR, 0026, err)
	}

	// 解析返回数据，提取交易
	if resp.Result == nil {
		return []entities.BaseDeposit {}, utils.LogMsgEx(utils.ERROR, "找不到指定块高的块：%d", height)
	}
	respData := resp.Result.(map[string]interface {})
	var txsObj interface {}
	var ok bool
	if txsObj, ok = respData["transactions"]; !ok {
		return nil, utils.LogIdxEx(utils.WARNING, 0001, nil)
	}

	// 定义提取属性的方法
	strProp := func(tx map[string]interface {}, key string) (string, error) {
		var itfc interface {}

		if itfc, ok = tx[key]; !ok {
			return "", utils.LogMsgEx(utils.ERROR, "交易未包含所需字段：%s", key)
		}

		return itfc.(string), nil
	}
	numProp := func(tx map[string]interface {}, key string) (*big.Int, error) {
		var err error
		var numTmp = big.NewInt(0)
		var strTmp string
		if strTmp, err = strProp(tx, key); err != nil {
			return numTmp, err
		}
		if strTmp[:2] == "0x" {
			strTmp = strTmp[2:]
		}

		if numTmp, ok = numTmp.SetString(strTmp, 16); !ok {
			return numTmp, utils.LogIdxEx(utils.ERROR, 29, strTmp)
		}
		return numTmp, nil
	}

	// 扫描所有交易
	txs := txsObj.([]interface {})
	deposits := []entities.BaseDeposit{}
	for i, tx := range txs {
		rawTx := tx.(map[string]interface {})
		deposit := entities.BaseDeposit{}
		if deposit.Address, err = strProp(rawTx, "to"); err != nil {
			deposit.Address = "create contract"
		}
		// 如果充值地址不属于钱包，跳过
		if !utils.StrArrayContains(addresses, deposit.Address) {
			continue
		}
		deposit.Asset	= "ETH"
		deposit.TxIndex	= i
		var heightBint *big.Int
		if heightBint, err = numProp(rawTx, "blockNumber"); err != nil {
			continue
		}
		deposit.Height = heightBint.Uint64()
		var timestampBint *big.Int
		if timestampBint, err = numProp(respData, "timestamp"); err != nil {
			continue
		}
		deposit.CreateTime = time.Unix(timestampBint.Int64(), 0)
		var valueBint *big.Int
		if valueBint, err = numProp(rawTx, "value"); err != nil {
			continue
		}
		var amountBflt = big.NewFloat(0)
		amountBflt.SetInt(valueBint.Div(valueBint, big.NewInt(int64(math.Pow10(rpc.decimal)))))
		deposit.Amount, _ = amountBflt.Float64()
		if deposit.TxHash, err = strProp(rawTx, "hash"); err != nil {
			continue
		}

		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

func (rpc *eth) GetCurrentHeight() (uint64, error) {
	var err error
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%d", rand.Intn(1000))
	if resp, err = rpc.sendRequest("eth_blockNumber", []interface{} {}, id); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 0026, err)
	}
	strHeight := resp.Result.(string)
	if strHeight[0:2] == "0x" {
		strHeight = strHeight[2:]
	}
	return strconv.ParseUint(strHeight, 16, 64)
}