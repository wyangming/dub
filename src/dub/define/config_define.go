package define

type DatabaseServerConfig struct {
	Addr string //绑定的ip地址及端口
	RegAddr string//注册服务器
}

type LogConfig struct {
	Level      string //日志级别
	DeviceName string //设备名字名字
	MaxSize    int    //按大小单位字节
}
type RegisterServerConfig struct {
	Addr string //绑定的ip地址及端口
}