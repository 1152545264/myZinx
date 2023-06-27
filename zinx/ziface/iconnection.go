package ziface

import "net"

// IConnection 定义连接模块的抽象
type IConnection interface {
	// Start 启动连接
	Start()

	// Stop 停止连接
	Stop()

	// GetTCPConnection 获取当前连接绑定的socket
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32

	// GetRemoteAddr 获取远程客户端的TCP状态 IP和Port
	GetRemoteAddr() net.Addr

	// SendMsg 发送数据，将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error

	//SetProperty 设置连接属性
	SetProperty(key string, value interface{})

	//GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)

	//RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
