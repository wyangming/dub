package secrpc

import (
	"dub/define"
	"dub/utils"
	"fmt"
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

	auth_all_args := &define.RpcDbAuthReqFindAllAuths{}
	auth_all_reply := &define.RpcDbAuthResFindAllAuths{}
	err = s.dbRpc.Call("AuthRpc.FindAllAuths", auth_all_args, auth_all_reply)
	if err != nil {
		s.log.Errorf("def_fun_sec_use.go LoginByLoginName method s.dbRpc.Call AuthRpc.FindAuthByRoleId err. %v\n", err)
		reply.Err = 2
		return nil
	}
	reply.Err = auth_all_reply.Err
	if reply.Err == 2 || auth_reply.Auths == nil || len(auth_all_reply.Auths) < 1 {
		return nil
	}

	use_entity_auth_map := make(map[uint]string)
	for _, auth := range auth_reply.Auths {
		if arg.ReqType == auth.AuthType {
			sec_auth := define.RpcSecUseResAuthModel{
				AuthId:              auth.AuthId,
				AuthPreId:           auth.AuthPreId,
				AuthName:            auth.AuthName,
				AuthMicroServerName: auth.AuthMicroServerName,
				AuthUrl:             auth.AuthUrl,
				AuthShowStatus:      auth.AuthShowStatus,
				AuthType:            auth.AuthType,
			}

			for _, enauth := range auth_all_reply.Auths {
				if auth.AuthId == enauth.AuthId && len(enauth.AuthUrl) > 0 {
					use_entity_auth_map[enauth.AuthId] = fmt.Sprintf("[%s]%s", enauth.AuthMicroServerName, enauth.AuthUrl)
				}

				if auth.AuthNeedUrl != nil && len(auth.AuthNeedUrl) > 0 {
					for _, needId := range auth.AuthNeedUrl {
						if needId == enauth.AuthId && len(enauth.AuthUrl) > 0 {
							use_entity_auth_map[enauth.AuthId] = fmt.Sprintf("[%s]%s", enauth.AuthMicroServerName, enauth.AuthUrl)
						}
					}
				}
			}

			reply.Auths = append(reply.Auths, sec_auth)
		}
	}

	//把用户的实体权限添加进去
	if len(use_entity_auth_map) > 0 {
		reply.PhyAuths = make([]string, 0, len(use_entity_auth_map))
		for _, val := range use_entity_auth_map {
			reply.PhyAuths = append(reply.PhyAuths, val)
		}
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
