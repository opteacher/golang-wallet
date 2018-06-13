package rpcs

import (
	"sync"
	"utils"
	"entities"
	"encoding/json"
	"log"
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

type Eth struct {
	sync.Once
	coinName string
	callUrl string
	decimal int
}

var __eth *Eth

func GetEth() *Eth {
	if __eth == nil {
		__eth = new(Eth)
		__eth.Once = sync.Once {}
		__eth.Once.Do(func() {
			__eth.create()
		})
	}
	return __eth
}

func (rpc *Eth) create() {
	setting := utils.GetConfig().GetCoinSettings()
	rpc.coinName 	= setting.Name
	rpc.callUrl		= setting.Url
	rpc.decimal		= setting.Decimal
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

func (rpc *Eth) sendRequest(method string, params []interface {}, id string) (EthSucceedResp, error)  {
	reqBody := RequestBody { method, params, id }
	reqStr, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Request body: %s", reqStr))

	reqBuf := bytes.NewBuffer([]byte(reqStr))
	res, err := http.Post(rpc.callUrl, "application/json", reqBuf)
	defer res.Body.Close()

	bodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Response body: %s", bodyStr))

	var resBody EthSucceedResp
	if err = json.Unmarshal(bodyStr, &resBody); err != nil {
		var resError EthFailedResp
		if err = json.Unmarshal(bodyStr, &resError); err != nil {
			log.Fatal(err)
		} else {
			log.Println(resError.Error)
			return resBody, errors.New(resError.Error.Message)
		}
	}
	return resBody, nil
}

func (rpc *Eth) GetTransactions(height uint) ([]entities.Deposit, error) {
	var err error
	var resp EthSucceedResp
	rand.Seed(time.Now().Unix())
	id := fmt.Sprintf("%d", rand.Intn(1000))
	params := []interface{} { "0x" + strconv.FormatUint(uint64(height), 16), true }
	if resp, err = rpc.sendRequest("eth_getBlockByNumber", params, id); err != nil {
		log.Println(err)
		return nil, err
	}

	respData := resp.Result.(map[string]interface {})
	var txsObj interface {}
	var ok bool
	if txsObj, ok = respData["transactions"]; !ok {
		err = errors.New("Error block chain response: no transactions in block")
		log.Println(err)
		return nil, err
	}

	strProp := func(tx map[string]interface {}, key string) (string, error) {
		var itfc interface {}

		if itfc, ok = tx[key]; !ok {
			msg := fmt.Sprintf("Transaction has no %s\n", key)
			log.Printf(msg)
			return "", errors.New(msg)
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
			err = errors.New(fmt.Sprintf("Cant convert to number: %s", strTmp))
			log.Println(err)
			return numTmp, err
		}
		return numTmp, nil
	}
	txs := txsObj.([]interface {})
	deposits := []entities.Deposit {}
	for i, tx := range txs {
		rawTx := tx.(map[string]interface {})
		deposit := entities.Deposit {}
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
		if deposit.Address, err = strProp(rawTx, "to"); err != nil {
			deposit.Address = "create contract"
		}
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}