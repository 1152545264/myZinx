package ziface

//定义一个服务器接口

type IServer interface {
	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	//AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)

	GetConnMgr() IConnManager

	//SetOnConnStart 注册OnConnStart钩子函数的方法
	SetOnConnStart(func(conn IConnection))

	//SetOnConnStop 注册OnConnStop钩子函数的方法
	SetOnConnStop(func(conn IConnection))

	//CallOnConnStart 调用OnConnStart钩子函数的方法
	CallOnConnStart(conn IConnection)

	//CallOnConnStop 调用OnConnStop钩子函数的方法
	CallOnConnStop(conn IConnection)
}
