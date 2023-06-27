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

// Test PreRouter
func (ph *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandler...")

	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping....\n"))
	if err != nil {
		fmt.Println("call back before ping error")

	}

}

// Test Handler
func (ph *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handler...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping...\n"))
	if err != nil {
		fmt.Println("call back ping..... ping error")

	}
}

// TestPostHandler
func (ph *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router Postandler...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping....\n"))
	if err != nil {
		fmt.Println("call back before ping error")

	}
}

func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[Zinx V0.4]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})

	//2启动Server
	s.Serve()
}
