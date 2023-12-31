package znet

import "zinx/ziface"

type Request struct {
	//已经和客户端建立好凡人连接
	conn ziface.IConnection

	//客户端请求的数据
	msg ziface.IMessage
}

// GetConnection 获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
