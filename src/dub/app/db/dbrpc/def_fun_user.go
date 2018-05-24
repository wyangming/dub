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

func (u *UserRpc) FindByName(loginName *string, reply *define.RpcDbUserResFindByName) error {
	dbAccounts := NewDbProxy().Get(DB_Accounts)
	row := dbAccounts.QueryRow("select userId,userName,loginName,loginPwd,userStatus,userAddTime,userAddId from dubuser where loginName=? limit 1", *loginName)

	var (
		userId, userAddId                             uint
		userStatus                                    uint8
		UserName, rowloginName, loginPwd, userAddTime string
	)
	err := row.Scan(&userId, &rowloginName, &rowloginName, &loginPwd, &userStatus, &userAddTime, &userAddId)

	switch {
	case err == sql.ErrNoRows:
		reply.Err = 1
	case err != nil:
		u.log.Errorf("def_fun_user.go FindByName method row.Scan err. \n", err)
	default:
		reply.UserId = userId
		reply.UserAddId = userAddId
		reply.UserStatus = userStatus
		reply.UserName = UserName
		reply.LoginName = rowloginName
		reply.LoginPwd = loginPwd

		if uatime, err := utils.TimeStrtoTime(userAddTime); err != nil {
			u.log.Errorf("def_fun_user.go FindByName method utils.TimeStrtoTime err. \n", err)
		} else {
			reply.UserAddTime = uatime
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

// insert into dubuserole (userId,roleId, useRoleStatus,urAddTime,urEditTime,urAddUserId)
// values(1,1,1,now(),now(),1) on duplicate key update urAddTime=now(),urEditTime=now(),useRoleStatus=1;
