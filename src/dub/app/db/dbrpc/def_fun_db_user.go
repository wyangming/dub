package dbrpc

import (
	"database/sql"
	"dub/define"
	"dub/utils"
)

//与用户有关的数据库操作
type UserRpc struct {
	log *utils.Logger
}

//根据用户查询用户表的信息
func (u *UserRpc) FindByName(loginName *string, reply *define.RpcDbUserResFindByName) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	row := dbAccounts.QueryRow("select userId,userName,loginName,loginPwd,userStatus,userAddTime,userAddId from dubuser where loginName=? limit 1", *loginName)

	var (
		userId, userAddId, userStatus                 sql.NullInt64
		UserName, rowloginName, loginPwd, userAddTime sql.NullString
	)
	err := row.Scan(&userId, &rowloginName, &rowloginName, &loginPwd, &userStatus, &userAddTime, &userAddId)

	switch {
	case err == sql.ErrNoRows:
		reply.Err = 1
	case err != nil:
		reply.Err = 2
		u.log.Errorf("def_fun_db_user.go FindByName method row.Scan err. %v\n", err)
	default:
		if userId.Valid {
			reply.UserId = uint(userId.Int64)
		}
		if userAddId.Valid {
			reply.UserAddId = uint(userAddId.Int64)
		}
		if userStatus.Valid {
			reply.UserStatus = uint8(userStatus.Int64)
		}
		if UserName.Valid {
			reply.UserName = UserName.String
		}
		if rowloginName.Valid {
			reply.LoginName = rowloginName.String
		}
		if loginPwd.Valid {
			reply.LoginPwd = loginPwd.String
		}

		if userAddTime.Valid {
			if uatime, err := utils.TimeStrtoTime(userAddTime.String); err != nil {
				u.log.Errorf("def_fun_db_user.go FindByName method utils.TimeStrtoTime err. %v\n", err)
			} else {
				reply.UserAddTime = uatime
			}
		}
	}
	return nil
}

var d_userRpc *UserRpc

func NewUserRpc() *UserRpc {
	if d_userRpc == nil {
		d_userRpc = &UserRpc{
			log: utils.NewLogger(),
		}
	}
	return d_userRpc
}
