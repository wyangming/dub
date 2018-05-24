package secrpc

import (
	"dub/utils"
	"fmt"
)

type SecRpc struct {
	log   *utils.Logger
	dbRpc *utils.RpcProxy
}

func (s *SecRpc) FindUserInfo(loginName, reply *string) error {
	*reply = fmt.Sprintf("%s-%s", *loginName, *reply)
	return nil
}

var _secRpc *SecRpc

func NewSecRpc(log *utils.Logger, dbRpc *utils.RpcProxy) *SecRpc {
	if _secRpc == nil {
		_secRpc = &SecRpc{
			log:   log,
			dbRpc: dbRpc,
		}
	}
	return _secRpc
}
