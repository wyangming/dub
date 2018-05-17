package common

type IConnector interface {
	//要连接服务器的地址
	Start(ip string) err
	//断开连接
	Stop() bool
	//发送命令
	Send(mainId uint16, subId uint16, data []byte) error
}
