package rpcs

import (
	"sync"
	"utils"
	"entities"
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"math/rand"
	"time"
	"math"
)

type btc struct {
	sync.Once
	coinName string
	callUrl string
	assSite string
	decimal int
	Stable int
	rpcUser string
	rpcPassword string
	account string
}

var __btc *btc

func (r *rpc) BTC() *btc {
	if __btc == nil {
		__btc = new(btc)
		__btc.Once = sync.Once {}
		__btc.Once.Do(func() {
			__btc.create()
		})
	}
	return __btc
}

func (rpc *btc) create() {
	setting := utils.GetConfig().GetCoinSettings()
	rpc.coinName 	= setting.Name
	rpc.callUrl		= setting.Url
	rpc.assSite		= setting.AssistSite
	rpc.decimal		= setting.Decimal
	rpc.Stable		= setting.Stable
	rpc.rpcUser		= setting.RPCUser
	rpc.rpcPassword	= setting.RPCPassword
	rpc.account		= setting.Deposit
}

type BtcResp struct {
	Error interface {}	`json:"error"`
	Id string			`json:"id"`
	Result interface{}	`json:"result"`
}

func (rpc *btc) sendRequest(method string, params []interface {}) (BtcResp, error) {
	id := fmt.Sprintf("%d", rand.Intn(1000))
	reqBody := RequestBody { method, params, id }
	reqStr, err := json.Marshal(reqBody)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 22, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Request body: %s", reqStr), nil)

	reqBuf := bytes.NewBuffer([]byte(reqStr))
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, rpc.callUrl, reqBuf); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 43, err))
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(rpc.rpcUser, rpc.rpcPassword)
	client := &http.Client {}
	res, err := client.Do(req)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 24, err))
	}
	defer res.Body.Close()

	bodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 25, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Response body: %s", bodyStr), nil)

	var resBody BtcResp
	if err = json.Unmarshal(bodyStr, &resBody); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 23, err))
	}
	if resBody.Result == nil {
		err = utils.LogIdxEx(utils.DEBUG, 40, nil)
	}
	if resBody.Error != nil {
		tmp := resBody.Error.(map[string]interface {})
		err = utils.LogIdxEx(utils.ERROR, 26, tmp["message"])
	}
	return resBody, err
}
func (rpc *btc) GetTransactions(height uint) ([]entities.Transaction, error) {
	var resp BtcResp
	var err error
	// 根据块高获取块哈希
	if resp, err = rpc.sendRequest("getblockhash", []interface {} { height }); err != nil {
		return []entities.Transaction {}, utils.LogIdxEx(utils.ERROR, 26, err)
	}
	blockHash := resp.Result.(string)

	// 根据块哈希获取块
	if resp, err = rpc.sendRequest("getblock", []interface {} { blockHash }); err != nil {
		return []entities.Transaction {}, utils.LogIdxEx(utils.ERROR, 26, err)
	}
	result := resp.Result.(map[string]interface {})

	var txs []entities.Transaction
	for _, txHash := range result["tx"].([]interface {}) {
		var txsTmp []entities.Transaction
		if txsTmp, err = rpc.GetTransaction(txHash.(string)); err != nil {
			utils.LogMsgEx(utils.ERROR, "找不到指定交易，交易HASH：%s，错误：%v", txHash, err)
			continue
		}
		txs = append(txs, txsTmp...)
	}
	return txs, nil
}
func (rpc *btc) GetCurrentHeight() (uint64, error) {
	var err error
	var resp BtcResp
	if resp, err = rpc.sendRequest("getblockcount", []interface {} {}); err != nil {
		utils.LogIdxEx(utils.ERROR, 26, err)
		return 0, err
	}
	return uint64(resp.Result.(float64)), nil
}
func (rpc *btc) GetDepositAmount() (map[string]float64, error) {
	if balance, err := rpc.GetBalance(rpc.account); err != nil {
		return map[string]float64 {}, err
	} else {
		return map[string]float64 { rpc.account: balance }, nil
	}
}
func (rpc *btc) GetBalance(address string) (float64, error) {
	var err error
	var resp BtcResp
	if resp, err = rpc.sendRequest("getbalance", []interface {} { address }); err != nil {
		return -1, utils.LogIdxEx(utils.ERROR, 26, err)
	}
	return resp.Result.(float64) * math.Pow10(-rpc.decimal), nil
}
func (rpc *btc) SendTransaction(from string, to string, amount float64, password string) (string, error) {
	params := []interface {} { from, to }
	params = append(params, math.Floor(amount * math.Pow10(rpc.decimal)))
	if resp, err := rpc.sendRequest("sendfrom", params); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 35, err)
	} else {
		return resp.Result.(string), nil
	}
}
func (rpc *btc) SendFrom(from string, amount float64) (string, error) {
	coinSet := utils.GetConfig().GetCoinSettings()
	return rpc.SendTransaction(rpc.account, coinSet.Collect, amount, rpc.account)
}
func (rpc *btc) SendTo(to string, amount float64) (string, error) {
	coinSet := utils.GetConfig().GetCoinSettings()
	return rpc.SendTransaction(coinSet.Withdraw, to, amount, coinSet.TradePassword)
}
func (rpc *btc) GetNewAddress() (string, error) {
	var resp BtcResp
	var err error
	if resp, err = rpc.sendRequest("getnewaddress", []interface {} { rpc.account }); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 36, err)
	}
	return resp.Result.(string), nil
}
func (rpc *btc) ValidAddress(address string) (bool, error) {
	if resp, err := rpc.sendRequest("validateaddress", []interface {} { address }); err != nil {
		return false, utils.LogMsgEx(utils.ERROR, "地址验证错误：%v", err)
	} else {
		result := utils.JsonObject {
			Data: resp.Result.(map[string]interface {}),
		}
		if !result.Contain("isvalid") {
			return false, nil
		}
		if tmp, err := result.Get("isvalid"); err != nil {
			return false, err
		} else {
			return tmp.(bool), nil
		}
	}
	return false, nil
}
func (rpc *btc) GetTransaction(txHash string) ([]entities.Transaction, error) {
	var err error
	var resp BtcResp
	if resp, err = rpc.sendRequest("getrawtransaction", []interface {} { txHash, 1 }); err != nil {
		return []entities.Transaction {}, utils.LogIdxEx(utils.ERROR, 37, err)
	}
	result := utils.JsonObject {
		Data: resp.Result.(map[string]interface {}),
	}
	blockTime, _ := result.Get("blocktime")
	blockHash, _ := result.Get("blockhash")

	var txs []entities.Transaction
	if !result.Contain("vout") {
		return nil, utils.LogMsgEx(utils.ERROR, "交易不含vout分量")
	}
	tmp, _ := result.Get("vout")
	for _, vout := range tmp.([]interface {}) {
		vos := utils.JsonObject { vout.(map[string]interface {}) }
		if tmp, err = vos.Get("scriptPubKey.type"); err != nil || tmp.(string) == "nulldata" {
			continue
		}
		if tmp, err = vos.Get("scriptPubKey.addresses"); err != nil {
			continue
		}
		var tx entities.Transaction
		tx.TxHash = txHash
		tx.CreateTime = time.Unix(int64(blockTime.(float64)), 0)
		if !vos.Contain("value") {
			continue
		}
		tmp, _ = vos.Get("value")
		tx.Amount = tmp.(float64) * math.Pow10(-rpc.decimal)
		tmp, _ = vos.Get("scriptPubKey.addresses")
		tx.To = tmp.([]interface {})[0].(string)
		//tx.From
		tx.BlockHash = blockHash.(string)
		if !vos.Contain("n") {
			continue
		}
		tmp, _ = vos.Get("n")
		tx.TxIndex = int(tmp.(float64))
		tx.Height, err = rpc.GetTxExistsHeight(txHash)
		if err != nil {
			continue
		}
		tx.Asset = rpc.coinName
		txs = append(txs, tx)
	}
	return txs, nil
}
func (rpc *btc) GetTxExistsHeight(txHash string) (uint64, error) {
	var err error
	var resp BtcResp
	if resp, err = rpc.sendRequest("gettransaction", []interface {} { txHash }); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 37, err)
	}
	result := utils.JsonObject {
		Data: resp.Result.(map[string]interface {}),
	}
	if !result.Contain("blockindex") {
		return 0, utils.LogMsgEx(utils.ERROR, "指定交易：%s中不存在块索引", txHash)
	}
	tmp, _ := result.Get("blockindex")
	return uint64(tmp.(float64)) + 1, nil
}
func (rpc *btc) EnableMining(enable bool, speed int) (bool, error) {
	return true, nil
}