package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"

	malt "github.com/taluos/Malt"
	"github.com/taluos/Malt/core/auth"
	consulRegistry "github.com/taluos/Malt/core/registry/consul"
	jwtauth "github.com/taluos/Malt/pkg/auth-jwt"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/log"
	restserver "github.com/taluos/Malt/server/rest"
	ginServer "github.com/taluos/Malt/server/rest/rest-gin"
)

func initAuthOperator() *auth.AuthOperator {
	userID := uuid.New().String()
	jwtInfo, err := JWT.NewJwtInfo(testPrivateKey, userID, "test:test", "admin", time.Minute*5)
	if err != nil {
		log.Fatalf("创建JWT信息失败: %v", err)
	}

	authOperator, _ := jwtauth.NewAuthenticator(jwtInfo.Keyfunc())

	jwtStratgy := auth.NewJWTStrategy(*jwtInfo, *authOperator, jwtInfo.Keyfunc())
	authOp := auth.AuthOperator{}
	authOp.SetStrategy(jwtStratgy)
	return &authOp
}

func main() {
	restServerSet := []restserver.Server{}

	restServerInstance := restserver.NewServer("gin",
		ginServer.WithPort(8080),
		ginServer.WithAuthOperator(initAuthOperator()),
		ginServer.WithMiddleware(gin.Recovery()),
	)

	restServerSet = append(restServerSet, restServerInstance)

	consulClient, err := api.NewClient(&api.Config{Address: "192.168.142.136:8500"})
	if err != nil {
		log.Fatalf("创建consul客户端失败: %v", err)
	}

	RegistyInstance := consulRegistry.New(
		consulClient,
		consulRegistry.WithHealthCheck(true),
		consulRegistry.WithHeartbeat(true),
		consulRegistry.WithHealthCheckInterval(10),
	)

	var App = malt.New(
		malt.WithId(uuid.New().String()),
		malt.WithName("Malt"),
		malt.WithTags([]string{"Rest:8080"}),
		malt.WithMetadata(map[string]string{"env": "dev", "Rest": "8080"}),
		malt.WithRegistrarTimeout(5*time.Second),
		malt.WithStopTimeout(5*time.Second),

		malt.WithRESTServer(restServerSet...),

		malt.WithRegistrar(RegistyInstance),
	)
	err = App.Run()
	if err != nil {
		log.Fatalf("server failed: %v", err)
		panic(err)
	}
}
