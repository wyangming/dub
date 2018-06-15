package dbrpc

import (
	"database/sql"
	"dub/define"
	"dub/utils"
)

//与用户有关的数据库操作
type RoleRpc struct {
	log *utils.Logger
}

//根据用户Id查询用户角色信息
func (r *RoleRpc) FindRoleByUseId(useId *uint, reply *define.RpcDbRoleResFindRoleByUseId) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	row := dbAccounts.QueryRow("select r.roleId,r.roleName,ur.urAddTime,ur.urEditTime,ur.urMain from dubuserole ur LEFT JOIN dubrole r on ur.roleId=r.roleId where r.roleStatus=1 and ur.useRoleStatus=1 and ur.urMain=1 and ur.userId=?", *useId)

	var (
		roleId, urMain                          sql.NullInt64
		roleName, urAddTime_str, urEditTime_str sql.NullString
	)
	err := row.Scan(&roleId, &roleName, &urAddTime_str, &urEditTime_str, &urMain)
	switch {
	case err == sql.ErrNoRows:
		reply.Err = 1
	case err != nil:
		reply.Err = 2
		r.log.Errorf("def_fun_db_auth.go FindRoleByUseId method row.Scan err. %v\n", err)
	default:
		if roleId.Valid {
			reply.RoleId = uint(roleId.Int64)
		}
		if roleName.Valid {
			reply.RoleName = roleName.String
		}
		if urMain.Valid && urMain.Int64 == 1 {
			reply.UrMain = true
		}

		if urAddTime_str.Valid {
			if urAddTime, err := utils.TimeStrtoTime(urAddTime_str.String); err != nil {
				r.log.Errorf("def_fun_db_auth.go FindRoleByUseId method utils.TimeStrtoTime err. %v\n", err)
			} else {
				reply.UrAddTime = urAddTime
			}
		}

		if urEditTime_str.Valid {
			if urEditTime, err := utils.TimeStrtoTime(urEditTime_str.String); err != nil {
				r.log.Errorf("def_fun_db_auth.go FindRoleByUseId method utils.TimeStrtoTime err. %v\n", err)
			} else {
				reply.UrEditTime = urEditTime
			}
		}
	}
	return nil
}

var d_roleRpc *RoleRpc

func NewRoleRpc() *RoleRpc {
	if d_roleRpc == nil {
		d_roleRpc = &RoleRpc{
			log: utils.NewLogger(),
		}
	}
	return d_roleRpc
}
