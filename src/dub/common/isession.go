package common

type CloseType int

var (
	CloseType_UpLayer    CloseType = 0 //上层主动
	CloseType_UnderLayer CloseType = 1 //下层主动
)

type ISession interface {
	// 发包
	Send(uint16, uint16, []byte) error
	// 断开
	Close(CloseType)
	// 标示ID
	ID() uint64
	//获取注册服务ip地址
	GetRemoteAddr() string
	//设置session数据
	SetSessionData(key string, val interface{})
	//得到session数据
	GetSessionData(key string) interface{}
}
