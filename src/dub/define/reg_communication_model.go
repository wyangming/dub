package define

type ModelRegReqServerType struct {
	Addr         string //服务地址
	ProtocolType uint8  //协议类型:0 rpc协议
	ServerType   uint8  //服务类型:0数据服务 1逻辑服务 2应用服务 3网关服务
	ServerName   string //服务名称
}

type ModelRegResServerType struct {
	Err uint8 //0成功 1请求类型格式错误
}
