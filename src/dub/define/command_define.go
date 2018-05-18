package define

//命令定义
const (
	Command_Reg = iota //注册服务器命令
)

//注册服务器子命令
const (
	CommandSub_Reg_Net        = iota //心跳包命令
	CommandSub_Reg_ServerType        //服务器类型命令
)
