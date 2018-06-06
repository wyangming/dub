package define

import (
	"time"
)

//用户名登录请求的信息
type RpcSecUseReqLoginByLoginName struct {
	LoginName, LoginPwd string //用户名密码
}

//用户名登录返回的信息
type RpcSecUseResLoginByLoginName struct {
	Err uint8 //0正常,1没有找到信息,2数据库连接错误,3用户名密码不正确

	UserId      uint
	UserName    string
	LoginName   string
	LoginPwd    string
	UserStatus  uint8
	UserAddTime time.Time
	UserAddId   uint

	RoleId                uint      //用户角色编号
	RoleName              string    //角色编号
	UrAddTime, UrEditTime time.Time //添加时间与最后修改时间
	UrMain                bool      //是否为主角色

	//权限信息
	Auths []RpcSecUseResAuthModel //所拥有的权限
}

//用户服务权限模型
type RpcSecUseResAuthModel struct {
	AuthId, AuthPreId                      uint   //权限编号 父级编号
	AuthName, AuthMicroServerName, AuthUrl string //权限名称 对应的微服务名称 权限链接
	AuthShowStatus, AuthType               uint8  //权限状态 权限类型
}
