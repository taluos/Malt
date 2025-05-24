package rest

var (
	ServerFactories = map[string]func(opts ...ServerOptions) Server{
		"gin":   newGinServer,
		"fiber": newFiberServer,
		// 可以注册其他服务器类型
	}
)
