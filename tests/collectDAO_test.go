package tests

import (
	"testing"
	"fmt"
	"dao"
)

func TestAddCollect(t *testing.T) {
	fmt.Println(dao.GetCollectDAO().AddSentCollect("0x1234", "ETH", "0xabcd", 120))
}