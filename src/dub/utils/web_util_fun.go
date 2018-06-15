package utils

import (
	"dub/define"
	"encoding/gob"
	"strconv"
)

//转换为uint类型
func WebTemplateC2Int(val interface{}) int {
	switch val.(type) {
	case uint8:
		res, _ := val.(uint8)
		return int(res)
	case uint16:
		res, _ := val.(uint16)
		return int(res)
	case uint32:
		res, _ := val.(uint32)
		return int(res)
	case uint64:
		res, _ := val.(uint64)
		return int(res)
	case uint:
		res, _ := val.(uint)
		return int(res)
	case int:
		res, _ := val.(int)
		return int(res)
	case string:
		str, _ := val.(string)
		res, _ := strconv.Atoi(str)
		return res
	default:
		return 0
	}
	return 0
}

//注册session存储的公共结构体
func RegSessionGobStruct() {
	gob.Register(define.RpcSecUseResLoginByLoginName{})
}
