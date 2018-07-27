package utils

import (
	"strings"
)

type JsonObject struct {
	Data map[string]interface {}
}

func (obj *JsonObject) Contain(key string) bool {
	tmp := obj.Data
	ks := strings.Split(key, ".")
	for i, k := range ks {
		if t, ok := tmp[k]; ok {
			if i == len(ks) - 1 {
				return true
			} else {
				switch t.(type) {
				case map[string]interface {}:
					tmp = t.(map[string]interface {})
				default:
					return false
				}
			}
			tmp = t.(map[string]interface {})
		} else {
			return false
		}
	}
	return true
}

func (obj *JsonObject) Get(key string) (interface {}, error) {
	tmp := obj.Data
	ks := strings.Split(key, ".")
	for i, k := range ks {
		if t, ok := tmp[k]; ok {
			if i == len(ks) - 1 {
				return t, nil
			} else {
				switch t.(type) {
				case map[string]interface {}:
					tmp = t.(map[string]interface {})
				default:
					return nil, LogMsgEx(ERROR, "不存在指定键值：%s(%s)", key, k)
				}
			}
		} else {
			return nil, LogMsgEx(ERROR, "不存在指定键值：%s(%s)", key, k)
		}
	}
	return nil, nil
}