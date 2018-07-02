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
	"dao"
)

type eth struct {
	sync.Once
	coinName string
	callUrl string
	decimal int
	Stable int
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
		return "", utils.LogMsgEx(utils.WARNING, "字段：%s为nil", key)
	}
	return itfc.(string), nil
}

func (rpc *eth) numProp(tx map[string]interface{}, key string) (*big.Int, error) {
	var err error
	var numTmp = big.NewInt(0)
	var strTmp string
	if strTmp, err = rpc.strProp(tx, key); err != nil {
		return numTmp, err
	}
	if strTmp[:2] == "0x" {
		strTmp = strTmp[2:]
	}
	var ok bool
	if numTmp, ok = numTmp.SetString(strTmp, 16); !ok {
		return numTmp, utils.LogIdxEx(utils.ERROR, 29, strTmp)
	}
	return numTmp, nil
}

func (rpc *eth) sendRequest(method string, params []interface {}, id string) (EthSucceedResp, error)  {
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

	var resBody EthSucceedResp
	if err = json.Unmarshal(bodyStr, &resBody); err != nil {
		var resError EthFailedResp
		if err = json.Unmarshal(bodyStr, &resError); err != nil {
			panic(utils.LogIdxEx(utils.ERROR, 23, err))
		} else {
			return resBody, utils.LogIdxEx(utils.ERROR, 26, errors.New(resError.Error.Message))
		}
	}
	if resBody.Result == nil {
		err = utils.LogIdxEx(utils.ERROR, 40, nil)
	}
	return resBody, err
}

func (rpc *eth) GetTransactions(height uint) ([]entities.Transaction, error) {
	var err error
	// 发送请求获取指定高度的块
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%d", rand.Intn(1000))
	params := []interface{} { "0x" + strconv.FormatUint(uint64(height), 16), true }
	if resp, err = rpc.sendRequest("eth_getBlockByNumber", params, id); err != nil {
		return nil, utils.LogIdxEx(utils.ERROR, 26, err)
	}

	// 解析返回数据，提取交易
	if resp.Result == nil {
		return []entities.Transaction {}, utils.LogMsgEx(utils.ERROR, "找不到指定块高的块：%d", height)
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
		var heightBint *big.Int
		if heightBint, err = rpc.numProp(rawTx, "blockNumber"); err != nil {
			continue
		}
		rpcTx.Height = heightBint.Uint64()
		var timestampBint *big.Int
		if timestampBint, err = rpc.numProp(respData, "timestamp"); err != nil {
			continue
		}
		rpcTx.CreateTime = time.Unix(timestampBint.Int64(), 0)
		var valueBint *big.Int
		if valueBint, err = rpc.numProp(rawTx, "value"); err != nil {
			continue
		}
		var amountBflt = big.NewFloat(0)
		amountBflt.SetInt(valueBint.Div(valueBint, big.NewInt(int64(math.Pow10(rpc.decimal)))))
		rpcTx.Amount, _ = amountBflt.Float64()
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
	id := fmt.Sprintf("%d", rand.Intn(1000))
	if resp, err = rpc.sendRequest("eth_blockNumber", []interface{} {}, id); err != nil {
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

func (rpc *eth) sendTransaction(from string, to string, amount float64) (string, error) {
	var err error
	coinSet := utils.GetConfig().GetCoinSettings()

	// 处理转账金额
	var amountFlt big.Float
	amountFlt.SetFloat64(amount)
	decimal := math.Pow10(utils.GetConfig().GetCoinSettings().Decimal)
	var decimalFlt big.Float
	decimalFlt.SetFloat64(decimal)
	var amountInt big.Int
	amountInt.SetString(amountFlt.Mul(&amountFlt, &decimalFlt).String(), 10)
	cvtAmount := fmt.Sprintf("0x%x", &amountInt)

	// 计算手续费数量
	var paramEstimateGas EstimateGasBody
	paramEstimateGas.From = from
	paramEstimateGas.To = to
	paramEstimateGas.Value = cvtAmount
	id := fmt.Sprintf("%d", rand.Intn(1000))
	var resp EthSucceedResp
	if resp, err = rpc.sendRequest("eth_estimateGas", []interface {} {
		paramEstimateGas,
	}, id); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 32, err)
	}
	gasNum := resp.Result

	// 解锁用户账户
	id = fmt.Sprintf("%d", rand.Intn(1000))
	if resp, err = rpc.sendRequest("personal_unlockAccount", []interface {} {
		from, coinSet.TradePassword, coinSet.UnlockDuration,
	}, id); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 33, err)
	}

	// 发送转账请求
	var paramTransaction TransactionBody
	paramTransaction.From = from
	paramTransaction.To = to
	paramTransaction.Value = cvtAmount
	paramTransaction.Gas = gasNum.(string)
	id = fmt.Sprintf("%d", rand.Intn(1000))
	if resp, err = rpc.sendRequest("eth_sendTransaction", []interface {} {
		paramTransaction,
	}, id); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 35, err)
	}
	if resp.Result == nil {
		return "", nil
	} else {
		return resp.Result.(string), nil
	}
}

func (rpc *eth) SendFrom(from string, amount float64) (string, error) {
	return rpc.sendTransaction(from, utils.GetConfig().GetCoinSettings().Collect, amount)
}

func (rpc *eth) SendTo(to string, amount float64) (string, error) {
	return rpc.sendTransaction(utils.GetConfig().GetCoinSettings().Withdraw, to, amount)
}

func (rpc *eth) GetBalance(address string) (float64, error) {
	var err error
	var resp EthSucceedResp
	id := fmt.Sprintf("%d", rand.Intn(1000))
	params := []interface {} { address, "latest" }
	if resp, err = rpc.sendRequest("eth_getBalance", params, id); err != nil {
		return -1, utils.LogIdxEx(utils.ERROR, 31, err)
	}

	result := resp.Result.(string)
	if result[0:2] == "0x" {
		result = result[2:]
	}
	var balanceInt big.Int
	balanceInt.SetString(result, 16)
	var balanceFlt big.Float
	decimal := big.NewInt(int64(math.Pow10(rpc.decimal)))
	balanceFlt.SetInt(balanceInt.Div(&balanceInt, decimal))
	ret, _ := balanceFlt.Float64()
	return ret, nil
}

func (rpc *eth) GetNewAddress() (string, error) {
	var err error
	var resp EthSucceedResp
	id := fmt.Sprintf("%d", rand.Intn(1000))
	params := []interface {} { utils.GetConfig().GetCoinSettings().TradePassword }
	if resp, err = rpc.sendRequest("personal_newAccount", params, id); err != nil {
		return "", utils.LogIdxEx(utils.ERROR, 36, err)
	}
	return resp.Result.(string), nil
}

func (rpc *eth) GetTransaction(txHash string) (entities.Transaction, error) {
	var err error
	var resp EthSucceedResp
	var tx entities.Transaction
	id := fmt.Sprintf("%d", rand.Intn(1000))
	if resp, err = rpc.sendRequest("eth_getTransactionByHash", []interface {} {
		txHash,
	}, id); err != nil {
		return tx, utils.LogIdxEx(utils.ERROR, 37, err)
	}
	result := resp.Result.(map[string]interface {})

	var numTmp *big.Int
	tx.TxHash = txHash
	tx.Asset = "ETH"
	if numTmp, err = rpc.numProp(result, "blockNumber"); err != nil {
		tx.Height = 0
	} else {
		tx.Height = numTmp.Uint64()
	}
	if numTmp, err = rpc.numProp(result, "transactionIndex"); err != nil {
		tx.TxIndex = -1
	} else {
		tx.TxIndex = int(numTmp.Int64())
	}
	tx.From = result["from"].(string)
	tx.To = result["to"].(string)
	tx.BlockHash = result["blockHash"].(string)
	if numTmp, err = rpc.numProp(result, "value"); err != nil {
		return tx, utils.LogIdxEx(utils.ERROR, 41, err)
	}
	fltTmp := big.NewFloat(0)
	fltTmp.SetInt(numTmp.Div(numTmp, big.NewInt(int64(math.Pow10(rpc.decimal)))))
	tx.Amount, _ = fltTmp.Float64()
	return tx, nil
}