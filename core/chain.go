package core

type HandlerFunc func(Ctx)

//type Handler func(Ctx)
//
//type Middleware func(Handler) Handler

/*
洋葱模型
传入业务函数Handler和中间件切片
倒序遍历中间件获取最终函数
*/
//func Chain(h Handler, mws ...Middleware) Handler {
//	for i := len(mws) - 1; i >= 0; i-- {
//		h = mws[i](h)
//	}
//	return h
//}
