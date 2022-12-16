package snet

import (
	"server/iface"
)

// 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写就好了
// 如果直接继承irouter类的话，需要将三个方法都进行实现，但是在实现中，不一定三个方法都需要进行实现，所以使用这种方法
type BaseRouter struct{}

// 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request iface.IRequest) {

}

// 在处理conn业务的主方法hook
func (br *BaseRouter) Handle(request iface.IRequest) {

}

// 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request iface.IRequest) {

}
