package znet

import "zinx/ziface"

// BaseRouter 实现router时，先嵌入这个baseRouter基类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct {
}

//BaseRouter的方法均为空是因为有的Router不希望有PreHand和PostHandle这两个方法

// PreHandle 在处理conn业务之前的钩子方法 HOOK
func (br *BaseRouter) PreHandle(request ziface.IRequest) {
}

// Handle 在处理conn业务的主方法hook
func (br *BaseRouter) Handle(request ziface.IRequest) {
}

// PostHandle 在处理conn业务之后的钩子方法HOOK
func (br *BaseRouter) PostHandle(request ziface.IRequest) {
}
