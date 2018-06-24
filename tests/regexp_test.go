package tests

import (
	"testing"
	"regexp"
	"fmt"
)

func TestRegExp(t *testing.T) {
	re := regexp.MustCompile("^/api/withdraw/([A-Z]{3,})/process$")
	fmt.Println(re.FindStringSubmatch("/api/withdraw/ETH/process"))
}
