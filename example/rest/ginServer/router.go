package restserver

import (
	restserver "github.com/taluos/Malt/server/rest"
)

func InitRouter(srv restserver.Server) {

	v1 := srv.Group("/v1")
	{
		userGroup := v1.Group("/user")
		{
			userController := NewUserServer()
			userGroup.Handle("GET", "/info", userController.GetUserInfo)
		}
	}
}
