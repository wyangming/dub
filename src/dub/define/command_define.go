package define

//注册服务器命令定义
const (
	CmdRegServer_Register = 1 //注册服务器的注册主命令
)

//注册服务器子命令
const (
	CmdSubRegServer_Register_Shut       = 0 //断开连接或者关闭服务命令
	CmdSubRegServer_Register_Reg        = 1 //向服务器注册命令
	CmdSubRegServer_Register_Reg_Inform = 2 //通知关联的服务器上线
)
