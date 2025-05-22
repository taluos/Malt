package rest

var (
	serverFactories = map[string]func(opts ...ServerOptions) Server{
		"gin": newGinServer,
		// 可以注册其他服务器类型
	}
)
