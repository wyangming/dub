package common

type IMoudle interface {
	//处理模块对应的消息
	OnMessage(id uint16, data []byte, ses ISession) bool
}
