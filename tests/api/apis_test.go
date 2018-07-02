package api

import (
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"apis"
	"github.com/stretchr/testify/assert"
	"reflect"
	"fmt"
)

func getFromServer(t *testing.T, path string) apis.RespVO {
	var resp *http.Response
	var err error
	if resp, err = http.Get(fmt.Sprintf("http://localhost:8037%s", path)); err != nil {
		t.Fatal(err)
	}

	var bodyBtary []byte
	if bodyBtary, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result apis.RespVO
	if err = json.Unmarshal(bodyBtary, &result); err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	return result
}

func TestGetHeightAPI(t *testing.T) {
	result := getFromServer(t, "/api/deposit/ETH/height")

	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "", result.Msg)
	assert.Equal(t, "float64", reflect.TypeOf(result.Data).Name())
}

func TestNewAddressAPI(t *testing.T) {
	result := getFromServer(t, "/api/deposit/ETH/address")

	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "", result.Msg)
	assert.Equal(t, "string", reflect.TypeOf(result.Data).Name())
}

func TestQueryProcess(t *testing.T) {
	result := getFromServer(t, "/api/common/ETH/process/0x6bccc827978af918acaf93e16f3887d9890e61e09dcbbe1a5b3676767cf9f232")

	assert.Equal(t, 200, result.Code)
	assert.Equal(t, "", result.Msg)
	data := result.Data.(map[string]interface {})
	assert.Equal(t, "0x6bccc827978af918acaf93e16f3887d9890e61e09dcbbe1a5b3676767cf9f232", data["tx_hash"])
}