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
	"strconv"
	"math/rand"
	"time"
	"math/big"
	"math"
	"dao"
)

type eth struct {
	sync.Once
	coinName string
	callUrl string
	decimal int
	Stable int
	isMining bool
}

var __eth *eth

func (r *rpc) ETH() *eth {
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
	rpc.isMining	= false
}

type EthSucceedResp struct {
	JsonRpc string		`json:"jsonrpc"`
	Id string			`json:"id"`
	Result interface{}	`json:"result"`
}

type EthFailedResp struct {
	JsonRpc string		`json:"jsonrpc"`
	Id string			`json:"id"`
	Error struct {
		Code int		`json:"code"`
		Message string	`json:"message"`
	}					`json:"error"`
}

type EstimateGasBody struct {
	From string		`json:"from"`
	To string		`json:"to"`
	Value string	`json:"value"`
}

type TransactionBody struct {
	EstimateGasBody
	Gas string		`json:"gas"`
}

func (rpc *eth) strProp(tx map[string]interface{}, key string) (string, error) {
	var itfc interface {}
	var ok bool
	if itfc, ok = tx[key]; !ok {
		return "", utils.LogMsgEx(utils.ERROR, "交易未包含所需字段：%s", key)
	}
	if itfc == nil {
		return "", nil// utils.LogMsgEx(utils.WARNING, "字段：%s为nil", key)
	}
	return itfc.(string), nil
}
func (rpc *eth) numProp(tx map[string]interface{}, key string) (*big.Float, error) {
	var err error
	var numTmp = big.NewFloat(0)
	var strTmp string
	if strTmp, err = rpc.strProp(tx, key); err != nil {
		return numTmp, err
	}
	if len(strTmp) == 0 {
		return numTmp, nil
	}
	if strTmp[:2] == "0x" {
		strTmp = strTmp[2:]
	}
	if numTmp, _, err = numTmp.Parse(strTmp, 16); err != nil {
		return numTmp, utils.LogIdxEx(utils.ERROR, 29, strTmp)
	}
	return numTmp, nil
}
func (rpc *eth) sendRequest(method string, params []interface {}) (EthSucceedResp, error)  {
	id := fmt.Sprintf("%d", rand.Intn(1000))
	reqBody := RequestBody { method, params, id }
	reqStr, err := json.Marshal(reqBody)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 22, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Request body: %s", reqStr), nil)

	reqBuf := bytes.NewBuffer([]byte(reqStr))
	res, err := http.Post(rpc.callUrl, "application/json", reqBuf)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 24, err))
	}
	defer res.Body.Close()

	bodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 25, err))
	}
	utils.LogMsgEx(utils.DEBUG, fmt.Sprintf("Response body: %s", bodyStr), nil)

	testBody := make(map[string]interface {})
	if err = json.Unmarshal(bodyStr, &testBody); err != nil {
		panic(utils.LogIdxEx(utils.ERROR, 23, err))
	}
	var resBody EthSucceedResp
	if _, ok := testBody["error"]; ok {
		var resError EthFailedResp
		if err = json.Unmarshal(bodyStr, &resError); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 23, err))
		} else {
			return resBody, utils.LogIdxEx(utils.ERROR, 26, resError.Error.Message)
		}
	} else {
		if err = json.Unmarshal(bodyStr, &resBody); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 23, err))
		}
	}
	if resBody.Result == nil {
		err = utils.LogIdxEx(utils.DEBUG, 40, nil)
	}
	return resBody, err
}
func (rpc *eth) GetTransactions(height uint) ([]entities.Transaction, error) {
	var err error
	// 发送请求获取指定高度的块
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	params := []interface{} { "0x" + strconv.FormatUint(uint64(height), 16), true }
	if resp, err = rpc.sendRequest("eth_getBlockByNumber", params); err != nil {
		return nil, utils.LogIdxEx(utils.ERROR, 26, err)
	}

	// 解析返回数据，提取交易
	if resp.Result == nil {
		return []entities.Transaction {}, utils.LogMsgEx(utils.DEBUG, "找不到指定块高的块：%d", height)
	}
	respData := resp.Result.(map[string]interface {})
	var txsObj interface {}
	var ok bool
	if txsObj, ok = respData["transactions"]; !ok {
		return nil, utils.LogIdxEx(utils.WARNING, 0001, nil)
	}

	// 扫描所有交易
	txs := txsObj.([]interface {})
	transactions := []entities.Transaction {}
	for i, tx := range txs {
		rawTx := tx.(map[string]interface {})
		rpcTx := entities.Transaction{}
		// 交易地址
		if rpcTx.From, err = rpc.strProp(rawTx, "from"); err != nil {
			rpcTx.From = ""
		}
		if rpcTx.To, err = rpc.strProp(rawTx, "to"); err != nil {
			rpcTx.To = ""
		}
		rpcTx.Asset = rpc.coinName
		rpcTx.TxIndex = i
		var heightBig *big.Float
		if heightBig, err = rpc.numProp(rawTx, "blockNumber"); err != nil {
			continue
		}
		rpcTx.Height, _ = heightBig.Uint64()
		var timestampBig *big.Float
		if timestampBig, err = rpc.numProp(respData, "timestamp"); err != nil {
			continue
		}
		timeInt64, _ := timestampBig.Int64()
		rpcTx.CreateTime = time.Unix(timeInt64, 0)
		var valueBig *big.Float
		if valueBig, err = rpc.numProp(rawTx, "value"); err != nil {
			continue
		}
		rpcTx.Amount, _ = valueBig.Mul(valueBig, big.NewFloat(math.Pow10(-rpc.decimal))).Float64()
		if rpcTx.TxHash, err = rpc.strProp(rawTx, "hash"); err != nil {
			continue
		}

		transactions = append(transactions, rpcTx)
	}
	return transactions, nil
}
func (rpc *eth) GetCurrentHeight() (uint64, error) {
	var err error
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	if resp, err = rpc.sendRequest("eth_blockNumber", []interface{} {}); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 26, err)
	}
	strHeight := resp.Result.(string)
	if strHeight[:2] == "0x" {
		strHeight = strHeight[2:]
	}
	return strconv.ParseUint(strHeight, 16, 64)
}
func (rpc *eth) GetDepositAmount() (map[string]float64, error) {
	var err error
	addrAmount := make(map[string]float64)
	coinCfg := utils.GetConfig().GetCoinSettings()

	var addresses []string
	addrDAO := dao.GetAddressDAO()
	if addresses, err = addrDAO.FindInuseByAsset(coinCfg.Name); err != nil {
		return addrAmount, utils.LogMsgEx(utils.ERROR, "获取充值地址失败：%v", err)
	}

	for _, addr := range addresses {
		var balance float64
		if balance, err = rpc.GetBalance(addr); err != nil {
			utils.LogMsgEx(utils.WARNING, "从%s获取余额失败", addr)
			continue
		}

		if balance < coinCfg.MinCollect {
			continue
		}

		addrAmount[addr] = balance
	}
	return addrAmount, nil
}
func (rpc *eth) SendTransaction(from string, to string, amount float64, password string) (string, error) {
	var err error
	coinSet := utils.GetConfig().GetCoinSettings()

	// 处理转账金额
	amountBig := big.NewFloat(amount)
	decimal := math.Pow10(rpc.decimal)
	amountBig.Mul(amountBig, big.NewFloat(decimal))
	amountFin := big.NewInt(0)
	amountBig.Int(amountFin)
	cvtAmount := fmt.Sprintf("0x%x", amountFin)

	// 计算手续费数量
	var paramEstimateGas EstimateGasBody
	paramEstimateGas.From = from
	paramEstimateGas.To = to
	paramEstimateGas.Value = cvtAmount
	var resp EthSucceedResp
	if resp, err = rpc.sendRequest("eth_estimateGas", []interface {} {
		paramEstimateGas,
	}); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 32, err)
	}
	gasNum := resp.Result

	// 解锁用户账户
	if resp, err = rpc.sendRequest("personal_unlockAccount", []interface {} {
		from, password, coinSet.UnlockDuration,
	}); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 33, err)
	}

	// 发送转账请求
	var paramTransaction TransactionBody
	paramTransaction.From = from
	paramTransaction.To = to
	paramTransaction.Value = cvtAmount
	paramTransaction.Gas = gasNum.(string)
	if resp, err = rpc.sendRequest("eth_sendTransaction", []interface {} {
		paramTransaction,
	}); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 35, err)
	}
	if resp.Result == nil {
		return "", nil
	} else {
		return resp.Result.(string), nil
	}
}
func (rpc *eth) SendFrom(from string, amount float64) (string, error) {
	coinSet := utils.GetConfig().GetCoinSettings()
	return rpc.SendTransaction(from, coinSet.Collect, amount, coinSet.TradePassword)
}
func (rpc *eth) SendTo(to string, amount float64) (string, error) {
	coinSet := utils.GetConfig().GetCoinSettings()
	return rpc.SendTransaction(coinSet.Withdraw, to, amount, coinSet.TradePassword)
}
func (rpc *eth) GetBalance(address string) (float64, error) {
	var err error
	var resp EthSucceedResp
	params := []interface {} { address, "latest" }
	if resp, err = rpc.sendRequest("eth_getBalance", params); err != nil {
		return -1, utils.LogIdxEx(utils.ERROR, 31, err)
	}
	if resp.Result == nil {
		return -1, nil
	}

	result := resp.Result.(string)
	if result[0:2] == "0x" {
		result = result[2:]
	}
	balanceBig := big.NewFloat(0)
	if balanceBig, _, err = balanceBig.Parse(result, 16); err != nil {
		return -1, utils.LogIdxEx(utils.ERROR, 41, err)
	}
	ret, _ := balanceBig.Mul(balanceBig, big.NewFloat(math.Pow10(-rpc.decimal))).Float64()
	return ret, nil
}
func (rpc *eth) GetNewAddress() (string, error) {
	var err error
	var resp EthSucceedResp
	params := []interface {} { utils.GetConfig().GetCoinSettings().TradePassword }
	if resp, err = rpc.sendRequest("personal_newAccount", params); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 36, err)
	}
	return resp.Result.(string), nil
}
func (rpc *eth) ValidAddress(address string) (bool, error) {
	if balance, err := rpc.GetBalance(address); err != nil || balance == -1 {
		return false, err
	} else {
		return true, nil
	}
}
func (rpc *eth) GetTransaction(txHash string) ([]entities.Transaction, error) {
	var err error
	var resp EthSucceedResp
	var tx entities.Transaction
	if resp, err = rpc.sendRequest("eth_getTransactionByHash", []interface {} {
		txHash,
	}); err != nil {
		return []entities.Transaction {}, utils.LogIdxEx(utils.ERROR, 37, err)
	}
	result := resp.Result.(map[string]interface {})

	var numTmp *big.Float
	tx.TxHash = txHash
	tx.Asset = "ETH"
	if numTmp, err = rpc.numProp(result, "blockNumber"); err != nil {
		tx.Height = 0
	} else {
		tx.Height, _ = numTmp.Uint64()
	}
	if numTmp, err = rpc.numProp(result, "transactionIndex"); err != nil {
		tx.TxIndex = -1
	} else {
		txIdxInt64, _ := numTmp.Int64()
		tx.TxIndex = int(txIdxInt64)
	}
	tx.From = result["from"].(string)
	tx.To = result["to"].(string)
	tx.BlockHash = result["blockHash"].(string)
	if numTmp, err = rpc.numProp(result, "value"); err != nil {
		return []entities.Transaction {}, utils.LogIdxEx(utils.ERROR, 41, err)
	}
	tx.Amount, _ = numTmp.Mul(numTmp, big.NewFloat(math.Pow10(-rpc.decimal))).Float64()
	return []entities.Transaction { tx }, nil
}
func (rpc *eth) GetTxExistsHeight(txHash string) (uint64, error) {
	var err error
	var resp EthSucceedResp
	if resp, err = rpc.sendRequest("eth_getTransactionByHash", []interface {} {
		txHash,
	}); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 37, err)
	}
	result := resp.Result.(map[string]interface {})

	var numTmp *big.Float
	var height uint64
	if numTmp, err = rpc.numProp(result, "blockNumber"); err != nil {
		return 0, utils.LogIdxEx(utils.ERROR, 37, err)
	} else {
		height, _ = numTmp.Uint64()
		return height, nil
	}
}
func (rpc *eth) EnableMining(enable bool, speed int) (bool, error)  {
	rpc.isMining = enable
	method := "miner_start"
	if !rpc.isMining {
		method = "miner_stop"
	}
	if _, err := rpc.sendRequest(method, []interface {} { speed }); err != nil {
		return false, utils.LogMsgEx(utils.ERROR, "调整挖矿状态失败：%v", err)
	}
	return true, nil
}
func (rpc *eth) IsMining() bool {
	return rpc.isMining
}