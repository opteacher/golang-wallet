package tests

import (
	"testing"
	"fmt"
	"dao"
)

func TestHeight(t *testing.T) {
	fmt.Println(dao.GetHeightDAO().GetHeight("ETH"))
}

func TestAddAsset(t *testing.T) {
	fmt.Println(dao.GetHeightDAO().ChkOrAddAsset("ETH"))
}

func TestUpdateHeight(t *testing.T) {
	fmt.Println(dao.GetHeightDAO().UpdateHeight("BTC", 10))
}