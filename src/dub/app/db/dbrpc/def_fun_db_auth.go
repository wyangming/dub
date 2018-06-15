package dbrpc

import (
	"database/sql"
	"dub/define"
	"dub/utils"
	"strconv"
	"strings"
)

//与权限有关的数据库操作
type AuthRpc struct {
	log *utils.Logger
}

func (a *AuthRpc) FindAllAuths(args *define.RpcDbAuthReqFindAllAuths, reply *define.RpcDbAuthResFindAllAuths) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	rows, err := dbAccounts.Query("select authId,authUrl,authMicroServerName from dubauthority")

	if err == sql.ErrNoRows {
		reply.Err = 1
		return nil
	}
	if err != nil && err != sql.ErrNoRows {
		reply.Err = 2
		a.log.Errorf("def_fun_db_auth.go FindAllAuths method dbAccounts.Query err.%v\n", err)
		return nil
	}

	defer rows.Close()
	reply.Auths = make([]*define.RpcDbAuthResAuthModel, 0)
	for rows.Next() {
		var (
			authId                       sql.NullInt64
			authUrl, authMicroServerName sql.NullString
		)
		err := rows.Scan(&authId, &authUrl, &authMicroServerName)
		if err != nil {
			reply.Err = 2
			a.log.Errorf("def_fun_db_auth.go FindAllAuths method rows.Scan err.\n", err)
		}
		auth := define.RpcDbAuthResAuthModel{}
		if authId.Valid {
			auth.AuthId = uint(authId.Int64)
		}
		if authUrl.Valid {
			auth.AuthUrl = authUrl.String
		}
		if authMicroServerName.Valid {
			auth.AuthMicroServerName = authMicroServerName.String
		}
		reply.Auths = append(reply.Auths, &auth)
	}
	return nil
}

//根据用户的编号查询相应的权限
func (a *AuthRpc) FindAuthByRoleId(roleId *uint, reply *define.RpcDbAuthResFindAuthByUseId) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	rows, err := dbAccounts.Query("select au.authId,au.authName,au.authPreId,au.authShowStatus,au.authMicroServerName"+
		",au.authType, au.authUrl,ra.roleAuthConf,au.authNeedUrlIds from dubauthority as au"+
		" LEFT JOIN dubroleauth as ra on au.authId=ra.authId where "+
		"ra.roleAuthStatus=1 and au.authStatus=1 and ra.roleId=? order by au.authAddTime asc", roleId)

	if err == sql.ErrNoRows {
		reply.Err = 1
		return nil
	}
	if err != nil && err != sql.ErrNoRows {
		reply.Err = 2
		a.log.Errorf("def_fun_db_auth.go FindAuthByUseId method dbAccounts.Query err.%v\n", err)
		return nil
	}

	defer rows.Close()
	reply.Auths = make([]*define.RpcDbAuthResAuthModel, 0)
	for rows.Next() {
		var (
			authId, authPreId, authShowStatus, authType, roleAuthConf sql.NullInt64
			authName, authMicroServerName, authUrl, authNeedUrl       sql.NullString
		)
		err := rows.Scan(&authId, &authName, &authPreId, &authShowStatus, &authMicroServerName, &authType, &authUrl, &roleAuthConf, &authNeedUrl)
		//err := rows.Scan(&auth.AuthId, &auth.AuthName, &auth.AuthPreId, &auth.AuthShowStatus, &auth.AuthMicroServerName, &auth.AuthType, &auth.AuthUrl)
		if err != nil {
			reply.Err = 2
			a.log.Errorf("def_fun_db_auth.go FindAuthByUseId method rows.Scan err.%v\n", err)
			continue
		}
		auth := define.RpcDbAuthResAuthModel{}
		if authId.Valid {
			auth.AuthId = uint(authId.Int64)
		}
		if authName.Valid {
			auth.AuthName = authName.String
		}
		if authPreId.Valid {
			auth.AuthPreId = uint(authPreId.Int64)
		}
		if authShowStatus.Valid {
			auth.AuthShowStatus = uint8(authShowStatus.Int64)
		}
		if authMicroServerName.Valid {
			auth.AuthMicroServerName = authMicroServerName.String
		}
		if authType.Valid {
			auth.AuthType = uint8(authType.Int64)
		}
		if authUrl.Valid {
			auth.AuthUrl = authUrl.String
		}
		if roleAuthConf.Valid {
			auth.RoleAuthConf = uint(roleAuthConf.Int64)
		}
		if authNeedUrl.Valid {
			ids := strings.Split(authNeedUrl.String, ",")
			if len(ids) > 0 {
				auth.AuthNeedUrl = make([]uint, 0)
				for _, id := range ids {
					if id != "" {
						if ind, err := strconv.Atoi(id); err != nil {
							a.log.Errorf("def_fun_db_auth.go FindAuthByRoleId method find data err au.authNeedUrl column is not int.%v\n", err)
						} else {
							auth.AuthNeedUrl = append(auth.AuthNeedUrl, uint(ind))
						}
					}
				}
			}
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
