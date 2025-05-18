package restserver

import (
	"Malt/pkg/log"

	"github.com/gin-gonic/gin"
)

type userServer struct {
}

func NewUserServer() *userServer {
	return &userServer{}
}

func (u *userServer) GetUserInfo(ctx *gin.Context) {
	log.Infof("get user info")
}
