package auth

import (
	// ginjwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"

	authJWT "github.com/taluos/Malt/pkg/auth-jwt"
	JWT "github.com/taluos/Malt/pkg/auth-jwt/JWT"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/pkg/log"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal"
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
func (j JWTStrategy) AuthFunc() fiber.Handler {
	if j.authenticator == nil {
		log.Errorf("jwt authenticator is nil")
	}
	return func(c fiber.Ctx) error {
		// 获取请求路径
		path := c.Route().Path // for example: /api/v1/hello
		// 获取 HTTP 方法
		method := c.Method() // for example: GET
		// 组合成 fullMethod
		fullMethod := method + ":" + path // for example: GET:/api/v1/hello
		authHeader := c.Get("Authorization")
		if err := j.authenticator.HTTPAuthenticate(authHeader, "", fullMethod, ""); err != nil {
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Token is not validable."), nil)
			c.Drop()
		}
		c.Next()
		return nil
	}
}
