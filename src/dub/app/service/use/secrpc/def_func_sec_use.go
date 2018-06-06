package secrpc

import (
	"dub/define"
	"dub/utils"
)

type SecUseRpc struct {
	log   *utils.Logger
	dbRpc *utils.RpcProxy
}

//用户名登录
func (s *SecUseRpc) LoginByLoginName(arg *define.RpcSecUseReqLoginByLoginName, reply *define.RpcSecUseResLoginByLoginName) error {
	//用户信息
	use_reply := define.RpcDbUserResFindByName{}
	err := s.dbRpc.Call("UserRpc.FindByName", &arg.LoginName, &use_reply)
	if err != nil {
		s.log.Errorf("def_fun_sec_use.go LoginByLoginName method s.dbRpc.Call UserRpc.FindByName err. %v\n", err)
		reply.Err = 2
		return nil
	}
	reply.Err = use_reply.Err
	if reply.Err == 2 {
		return nil
	}
	//用户名密码是否正确
	if use_reply.LoginName != arg.LoginName || use_reply.LoginPwd != arg.LoginPwd {
		reply.Err = 3
		return nil
	}
	reply.UserId = use_reply.UserId
	reply.UserName = use_reply.UserName
	reply.LoginName = use_reply.LoginName
	reply.LoginPwd = use_reply.LoginPwd
	reply.UserStatus = use_reply.UserStatus
	reply.UserAddTime = use_reply.UserAddTime
	reply.UserAddId = use_reply.UserAddId

	//角色信息
	role_reply := define.RpcDbRoleResFindRoleByUseId{}
	err = s.dbRpc.Call("RoleRpc.FindRoleByUseId", &reply.UserId, &role_reply)
	if err != nil {
		s.log.Errorf("def_fun_sec_use.go LoginByLoginName method s.dbRpc.Call RoleRpc.FindRoleByUseId err. %v\n", err)
		reply.Err = 2
		return nil
	}
	reply.Err = role_reply.Err
	if reply.Err == 2 {
		return nil
	}
	reply.RoleId = role_reply.RoleId
	reply.RoleName = role_reply.RoleName
	reply.UrAddTime = role_reply.UrAddTime
	reply.UrEditTime = role_reply.UrEditTime
	reply.UrMain = role_reply.UrMain

	//权限信息
	auth_reply := define.RpcDbAuthResFindAuthByUseId{}
	err = s.dbRpc.Call("AuthRpc.FindAuthByRoleId", &reply.RoleId, &auth_reply)
	if err != nil {
		s.log.Errorf("def_fun_sec_use.go LoginByLoginName method s.dbRpc.Call AuthRpc.FindAuthByRoleId err. %v\n", err)
		reply.Err = 2
		return nil
	}
	reply.Err = role_reply.Err
	if reply.Err == 2 || auth_reply.Auths == nil || len(auth_reply.Auths) < 1 {
		return nil
	}
	for _, auth := range auth_reply.Auths {
		sec_auth := define.RpcSecUseResAuthModel{
			AuthId:              auth.AuthId,
			AuthPreId:           auth.AuthPreId,
			AuthName:            auth.AuthName,
			AuthMicroServerName: auth.AuthMicroServerName,
			AuthUrl:             auth.AuthUrl,
			AuthShowStatus:      auth.AuthShowStatus,
			AuthType:            auth.AuthType,
		}
		reply.Auths = append(reply.Auths, sec_auth)
	}

	return nil
}

var _secRpc *SecUseRpc

func NewSecUseRpc(log *utils.Logger, dbRpc *utils.RpcProxy) *SecUseRpc {
	if _secRpc == nil {
		_secRpc = &SecUseRpc{
			log:   log,
			dbRpc: dbRpc,
		}
	}
	return _secRpc
}
