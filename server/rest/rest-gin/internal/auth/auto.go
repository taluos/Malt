package auth

import (
	"strings"

	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	internal "github.com/taluos/Malt/server/rest/rest-gin/internal"

	"github.com/gin-gonic/gin"
)

// AutoStrategy defines authentication strategy which can automatically choose between Basic and Bearer
// according `Authorization` header.
type AutoStrategy struct {
	basic BasicStrategy
	jwt   JWTStrategy
}

var _ AuthStrategy = &AutoStrategy{}

// NewAutoStrategy create auto strategy with basic strategy and jwt strategy.
func NewAutoStrategy(basic BasicStrategy, jwt JWTStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

// AuthFunc defines auto strategy as the gin authentication middleware.
func (a *AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := AuthOperator{}
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			internal.WriteResponse(
				c,
				errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(&a.basic)
		case "Bearer":
			operator.SetStrategy(&a.jwt)
			// a.JWT.MiddlewareFunc()(c)
		default:
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil)
			c.Abort()

			return
		}

		operator.AuthFunc()(c)

		c.Next()
	}
}
