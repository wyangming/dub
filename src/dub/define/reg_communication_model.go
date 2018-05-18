package define

type ModelRegReqServerType struct {
	Addr         string //服务地址
	ProtocolType uint8  //协议类型:0 rpc协议
	ServerType   uint8  //server type:0 db server
	ServerName   string //服务名称
}
