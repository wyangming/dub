与服务器通信命令规划：
0网络命令：
    1：心中跳包

1注册命令：
    0：断开连接或者关闭服务命令
    1：向服务器注册命令
    2：通知关联的服务器上线

注册模块：
注册服务器 模型完成
数据服务器 模型完成
逻辑服务器 模型完成
应用服务器 模型完成
网关服务器 网关完成


beego的session如果存储在文件里，它对文件不是独占的。做微服务时可以把本地的多个session存储在同一个目录下以实现session共享。但需要通过gob注册存储在session里的结构体