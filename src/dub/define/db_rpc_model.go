package define

import (
	"time"
)

//查询用户
type RpcDbUserResFindByName struct {
	Err         uint8 //0正常，1没有找到
	UserId      uint
	UserName    string
	LoginName   string
	LoginPwd    string
	UserStatus  uint8
	UserAddTime time.Time
	UserAddId   uint
}
