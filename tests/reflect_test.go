package tests

import (
	"testing"
	"fmt"
	"reflect"
)

func TestUseNameGetValue(t *testing.T) {
	fmt.Println(reflect.ValueOf("RFC3339").String())
}