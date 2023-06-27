package main

import "zinx/znet"

// 基于zinx框架来开发的 服务器端应用程序

func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.1]")

	//2启动Server
	s.Serve()
}
