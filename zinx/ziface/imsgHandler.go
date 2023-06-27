package ziface

// IMsgHandle  消息管理抽象层
type IMsgHandle interface {

	//DoMsgHandler 调度执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	//AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	//StartWorkerPool 启动Worker工作池
	StartWorkerPool()

	//SendMessageToQueue 将消息发送给任务队列进行处理
	SendMessageToQueue(request IRequest)
}
