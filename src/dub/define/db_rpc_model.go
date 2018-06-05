package define

import (
	"time"
)

//查询用户
type RpcDbUserResFindByName struct {
	Err         uint8 //0正常,1没有找到,2数据库连接错误
	UserId      uint
	UserName    string
	LoginName   string
	LoginPwd    string
	UserStatus  uint8
	UserAddTime time.Time
	UserAddId   uint
}

//用户的角色信息
type RpcDbRoleResFindRoleByUseId struct {
	Err                   uint8     //0正常,1没有找到,2数据库连接错误
	RoleId                uint      //用户角色编号
	RoleName              string    //角色编号
	UrAddTime, UrEditTime time.Time //添加时间与最后修改时间
	UrMain                bool      //是否为主角色
}

//查询用户权限的所有权限
type RpcDbAuthResFindAuthByUseId struct {
	Err   uint8 //0正常,1没有找到,2数据库连接错误
	Auths []*RpcDbAuthResAuthModel
}

//用户的所有权限
type RpcDbAuthResAuthModel struct {
	AuthId, AuthPreId                      uint   //权限编号 父级编号
	AuthName, AuthMicroServerName, AuthUrl string //权限名称 对应的微服务名称 权限链接
	AuthShowStatus, AuthType               uint8  //权限状态 权限类型
}
