package define

type ModelRegReqServerType struct {
	Addr         string //服务地址
	ProtocolType uint8  //协议类型:0 rpc协议
	//服务类型:0数据服务 1逻辑服务 2web应用服务 3网关服务 4是api应用服务
	//0->1 1->2 2->3 4->3
	//3网关服务现分为:web网关 api网关
	//根据不同的网关需要接入不再的应用服务，如api网关接入api应用服务,web网关接入web的应用服务
	ServerType uint8
	ServerName string //服务名称
	ProxyUrl   string //被网关代理的路径 目前只有2 4类型才有值，如果没有默认路径是/
}

type ModelRegResServerType struct {
	Err uint8 //0成功 1请求类型格式错误
}
