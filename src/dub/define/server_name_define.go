package define

//除了注册服务器没有在这里定义其他的都要在这里定义

//数据服务 type 0
const (
	ServerNameDB_DbServer = "db_DbServer"
)

//逻辑服务 type 1
const (
	ServerNameService_UseServer = "Service_UseServer"
)

//web服务器 type 2
const (
	ServerNameWeb_ManLobbyServer = "Service_ManLobbyServer"
)

//网关服务器 type 3
const (
	ServerNameGate_WebServer = "Gate_WebServer"
)
