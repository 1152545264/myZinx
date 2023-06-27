package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// 基于zinx框架来开发的 服务器端应用程序

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handler
func (ph *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handler...")
	//读取客户端的数据 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[Zinx V0.5]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})

	//2启动Server
	s.Serve()
}
