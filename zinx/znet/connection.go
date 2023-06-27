package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

// Connection 连接模块
type Connection struct {
	//当前连接的socket
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前连接已经退出的/停止 channel （由Reader告知Writer退出）
	ExitChan chan bool

	//无缓冲管道，用于读写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理业务API
	MsgHandler ziface.IMsgHandle

	//当前connection属于哪个server
	TcpServer ziface.IServer

	//连接属性集合
	property map[string]interface{}
	//保护连接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		TcpServer:  server,
		property:   make(map[string]interface{}),
	}
	//将conn加入到ConnManager
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// StartReader 连接的读业务方法，对应于服务器解包过程
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit], remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()
	for {

		//创建一个拆包解包对象
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包 得到msgID 和 msgDataLen放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}

		//根据dataLen 再次读取data ，放在msg.Data中
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			msg.SetMsgData(data)
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMessageToQueue(&req)
		} else {

			//从路由中，找到注册绑定的Conn对应的router调用
			//根据绑定好的MsgID找到对应的处理api业务进行执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}

}

// StartWriter 写消息Goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.GetRemoteAddr().String(), "[conn Writer exit!]")

	//不断的阻塞等待channel的消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也需要退出
			return
		}
	}
}

// Start 启动连接
func (c *Connection) Start() {
	fmt.Println("Conn Start()....., ConnID=", c.ConnID)

	//下面实现了读写分离的业务
	//启动从当前连接的读数据的业务
	go c.StartReader()

	//启动从当前连接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的创建连接之后需要调用的处理业务 执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)

}

// Stop 停止连接
func (c *Connection) Stop() {
	fmt.Println("Conn stop..., ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//销毁连接之前调用开发者注册的hook 需要执行的业务hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)

	//将当前连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)
}

// GetTCPConnection 获取当前连接绑定的socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetRemoteAddr 获取远程客户端的TCP状态 IP和Port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 发送数据，将数据发送给远程的客户端。对应于服务器封包过程
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//将data进行封包 MsgDataLen | MsgID | MsgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessagePackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg ")
	}

	c.msgChan <- binaryMsg

	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
