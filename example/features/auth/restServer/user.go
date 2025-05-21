package restserver

import "github.com/gin-gonic/gin"

type userServer struct {
}

func NewUserServer() *userServer {
	return &userServer{}
}

func (u *userServer) GetUserInfo(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"msg": ""})
}
