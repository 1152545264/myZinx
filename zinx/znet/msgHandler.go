package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

//消息处理模块的具体实现

type MessageHandler struct {

	//存放每个MsgID对应的处理方法
	APIS map[uint32]ziface.IRouter

	//负责Worker 取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的Worker数量
	WorkerPoolSize uint32
}

// NewMessageHandler 初始化创建MessageHandler方法
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		APIS:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置文件中获取,也可以在配置文件让用户进行设置
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度执行对应的Router消息处理方法
func (m *MessageHandler) DoMsgHandler(request ziface.IRequest) {
	//从request中找到msgID
	handler, ok := m.APIS[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! need to register")
		return
	}

	//2 找到msgID，调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (m *MessageHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前message绑定的API处理方法是否已经存在
	if _, ok := m.APIS[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	//2 添加message与API的绑定关系
	m.APIS[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个Worker工作池
func (m *MessageHandler) StartWorkerPool() {
	//根据WorkerPoolSize 分别开启Worker ，每个Worker用一个goroutine承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个Worker被启动
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前worker 阻塞等待消息从channel中传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (m *MessageHandler) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID = ", workID, " is started ....")

	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request,指定当前request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// SendMessageToQueue  将消息交给TaskQueue，由Worker进行处理
func (m *MessageHandler) SendMessageToQueue(request ziface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的ConnID来进行分配
	//基本的平均分配法则
	workID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID=", request.GetMsgID(),
		" to workerID = ", workID)

	//2 将消息发送给对应的worker对应的TaskQueue即可
	m.TaskQueue[workID] <- request
}
