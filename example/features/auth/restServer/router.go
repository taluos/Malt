package restserver

import (
	restserver "github.com/taluos/Malt/server/rest/Server"
)

func InitRouter(g *restserver.Server) {
	v1 := g.Group("/v1")
	{
		userGroup := v1.Group("/user")
		{
			userController := NewUserServer()
			userGroup.GET("info", userController.GetUserInfo)
		}
	}

}
