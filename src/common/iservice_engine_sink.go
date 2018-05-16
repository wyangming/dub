package common

//服务引擎信息
type IServiceEnginSink interface {
	//起动服务
	Serve() error
	//停止服务
	Stop()
	////与其他服务器连接
	//OnEventTCPSocketLink(IConnector) error
	////关闭与其他服务器连接
	//OnEventTCPSocketShut(IConnector) bool
	////读取其他服务器发送的信息
	//OnEventTCPSocketRead(uint16, uint16, []byte, IConnector) bool
	////其他服务器连接
	//OnEventTCPNetworkLink(ISession) bool
	////其他服务器与本机断开
	//OnEventTCPNetworkShut(ISession) bool
	////处理其他服务器与本机的请求
	//OnEventTCPNetworkRead(uint16, uint16, []byte, ISession) bool
}