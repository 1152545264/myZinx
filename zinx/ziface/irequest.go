package ziface

//IRequest接口:实际上是把客户端请求的连接信息和连接数据包装在了一个request请求中

type IRequest interface {
	//GetConnection 获取当前连接
	GetConnection() IConnection

	//GetData 获取请求的消息数据
	GetData() []byte

	//GetMsgID 获取请求消息的ID
	GetMsgID() uint32
}
