package define

//dbserver
type DatabaseServerConfig struct {
	Addr    string //绑定的ip地址及端口
	RegAddr string //注册服务器
}

//日志
type LogConfig struct {
	Level      string //日志级别
	DeviceName string //设备名字名字
	MaxSize    int    //按大小单位字节
}

//注册服务器
type RegisterServerConfig struct {
	Addr string //绑定的ip地址及端口
}

//用户服务
type ServiceUseServerConfig struct {
	Addr    string //绑定的ip地址及端口
	RegAddr string //注册服务器
}

//用户中心
type WebUserCenterServerConfig struct {
	Addr          string //web服务的ip地址及端口
	WebWiew       string //web页面的路径地址
	WebStaticPath string //web页面路径的静态地址
	WebStaticUrl  string //web页面的静态url
	RunMode       string //运行模式

	ProxyUrl string //被gate代理的url
	RegAddr  string //注册服务器
}

//注册服务器
type GateWebServerConfig struct {
	Addr    string //绑定的ip地址及端口
	RegAddr string //注册服务器
}
