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

// Handle Test Handler
func (ph *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handler...")
	//读取客户端的数据 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}

}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// DoConnectionBegin 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=======> DoConnectionBegin is Called......")
	if err := conn.SendMsg(202, []byte("Do Connection Begin")); err != nil {
		fmt.Println(err)
	}
}

// DoConnectionEnd 断开连接之前需要执行的函数
func DoConnectionEnd(conn ziface.IConnection) {
	fmt.Println("=======> DoConnectionEnd is Called......")
	fmt.Println("Conn ID = ", conn.GetConnID())
}

// Handle Test Handler
func (ph *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handler...")
	//读取客户端的数据 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[Zinx V0.9]")

	//注册连接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionEnd)

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//2启动Server
	s.Serve()
}
