package dbrpc

import (
	"database/sql"
	"dub/define"
	"dub/utils"
)

//与权限有关的数据库操作
type AuthRpc struct {
	log *utils.Logger
}

//根据用户的编号查询相应的权限
func (a *AuthRpc) FindAuthByRoleId(roleId *uint, reply *define.RpcDbAuthResFindAuthByUseId) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	rows, err := dbAccounts.Query("select au.authId,au.authName,au.authPreId,au.authShowStatus,au.authMicroServerName"+
		",au.authType, au.authUrl from dubauthority as au"+
		" LEFT JOIN dubroleauth as ra on au.authId=ra.authId where "+
		"ra.roleAuthStatus=1 and au.authStatus=1 and ra.roleId=? order by au.authAddTime asc", roleId)

	if err == sql.ErrNoRows {
		reply.Err = 1
		return nil
	}
	if err != nil {
		reply.Err = 2
		a.log.Errorf("def_fun_user.go FindAuthByUseId method dbAccounts.Query err.%v\n", err)
		return nil
	}

	defer rows.Close()
	reply.Auths = make([]*define.RpcDbAuthResAuthModel, 0)
	for rows.Next() {
		auth := define.RpcDbAuthResAuthModel{}
		err := rows.Scan(&auth.AuthId, &auth.AuthName, &auth.AuthPreId, &auth.AuthShowStatus, &auth.AuthMicroServerName, &auth.AuthType, &auth.AuthUrl)
		if err != nil {
			reply.Err = 2
			a.log.Errorf("def_fun_user.go FindAuthByUseId method rows.Scan err.%v\n", err)
			continue
		}

		reply.Auths = append(reply.Auths, &auth)
	}
	return nil
}

var d_authRpc *AuthRpc

func NewAuthRpc() *AuthRpc {
	if d_authRpc == nil {
		d_authRpc = &AuthRpc{
			log: utils.NewLogger(),
		}
	}
	return d_authRpc
}
