package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现， 定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的IP地址版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务API路由
	MsgHandler ziface.IMsgHandle

	//server 的连接管理器
	ConnMgr ziface.IConnManager

	//该server创建连接之后自动调用的Hook函数——OnConnStart
	OnCOnnStart func(conn ziface.IConnection)

	//该server销毁连接之前自动调用的Hook函数——OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listener at IP : %s, Port : %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version %s, MaxConn %d , MaxPackageSize: %d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启消息队列以及worker工作池
		s.MsgHandler.StartWorkerPool()

		//1 获取一个tcp的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err", err)
			return
		}
		fmt.Println("start Zinx Server success, ", s.Name, " success, Listening...")
		var cid uint32
		cid = 0

		//3 阻塞等待客户端连接， 处理客户端的连接业务（读写）
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}

			//设置最大连接数
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 給客戶端响应一个超出最大连接的错误包
				fmt.Print("Too many Connections, MaxConn = ", utils.GlobalObject.MaxConn)

				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn进行绑定得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	// 将一些服务器资源状态或者已经申请的连接信息进行停止或者回收
	fmt.Println("[STOP] Zinx server ", s.Name)
	s.ConnMgr.ClearConn()
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hook func(conn ziface.IConnection)) {
	s.OnCOnnStart = hook
}

// SetOnConnStop 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hook func(conn ziface.IConnection)) {
	s.OnConnStop = hook
}

// CallOnConnStart 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnCOnnStart != nil {
		fmt.Println("----> Call OnConnStart() ....")
		s.OnCOnnStart(conn)
	}

}

// CallOnConnStop 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println(" ------>Call OnConnStop() .....")
		s.OnConnStop(conn)
	}
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMessageHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
