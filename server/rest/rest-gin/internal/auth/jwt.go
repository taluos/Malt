package auth

import (
	// ginjwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	authJWT "github.com/taluos/Malt/pkg/auth-jwt"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/pkg/log"
	internal "github.com/taluos/Malt/server/rest/rest-gin/internal"
)

// AuthzAudience defines the value of jwt audience field.
const AuthzAudience = "Malt"

// JWTStrategy defines jwt bearer authentication strategy.
type JWTStrategy struct {
	JWT.JwtInfo
	authenticator *authJWT.Authenticator
	keyFunc       jwt.Keyfunc
}

var _ AuthStrategy = &JWTStrategy{}

func NewJWTStrategy(jwtInfo JWT.JwtInfo, auth authJWT.Authenticator, keyFunc jwt.Keyfunc) JWTStrategy {
	return JWTStrategy{jwtInfo, &auth, keyFunc}
}

// AuthFunc defines jwt bearer strategy as the gin authentication middleware.
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	if j.authenticator == nil {
		log.Errorf("jwt authenticator is nil")
	}
	return func(c *gin.Context) {
		// 获取请求路径
		path := c.Request.URL.Path
		// 获取 HTTP 方法
		method := c.Request.Method
		// 组合成 fullMethod
		fullMethod := method + ":" + path
		authHeader := c.GetHeader("Authorization")
		if err := j.authenticator.HTTPAuthenticate(authHeader, "", fullMethod, ""); err != nil {
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Token is not validable."), nil)
			c.Abort()
		}
		c.Next()
	}
}
