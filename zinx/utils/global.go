package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

// 存储一切有关Zinx框架的全局参数供其他模块使用
type GlobalOgj struct {
	//Server
	TCPServer ziface.IServer //zinx全局的Server对象
	Host      string         //当前服务器监听的IP地址
	TCPPort   int            //当前服务器监听的端口号
	Name      string         //	当前服务器的名称

	//version
	Version          string //当前zinx的版本号
	MaxConn          int    //当前服务器允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //zinx框架允许用户最多开辟多少个worker(限定条件)
}

// GlobalObject 定义一个全局的对外GlobalObj
var GlobalObject *GlobalOgj

// Reload 从zinx.json去加载用于自定义的参数
func (g *GlobalOgj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法
func init() {
	//如果配置文件没有加载的默认值
	GlobalObject = &GlobalOgj{
		Name:             "ZinxServerApp",
		Version:          "V0.9",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   100,  //整个Worker工作池队列的个数
		MaxWorkerTaskLen: 1024, //每个Worker对应的消息队列的任务数量最大值
	}

	//应该从配置文件中加载一些用户自定义的参数
	GlobalObject.Reload()
}
