package secrec

import (
	"dub/utils"
)

const (
	ConstServiceUseRpc = iota
)

//rpc服务信息
var serviceRpc map[uint8]*utils.RpcProxy

func AddSecRpc(rpcType uint8, rpc *utils.RpcProxy) {
	if serviceRpc == nil {
		serviceRpc = make(map[uint8]*utils.RpcProxy)
	}
	serviceRpc[rpcType] = rpc
}

func GetRpc(rpcType uint8) *utils.RpcProxy {
	if serviceRpc == nil {
		return nil
	}
	rpc, ok := serviceRpc[rpcType]
	if ok {
		return rpc
	}
	return nil
}

var proxyServerNames map[string]string

//服务器代理信息
func UpdateProxyInfo(proxyUrls, serverNames []string) {
	if len(proxyUrls) < 1 || len(serverNames) < 1 {
		return
	}

	if proxyServerNames == nil {
		proxyServerNames = make(map[string]string)
	}

	for i := 0; i < len(proxyUrls); i++ {
		proxyServerNames[serverNames[i]] = proxyUrls[i]
	}
}

func GetProxyByServerName(serverName string) string {
	if proxyServerNames == nil {
		return ""
	}
	proxyUrl, ok := proxyServerNames[serverName]
	if ok {
		return proxyUrl
	}
	return ""
}
