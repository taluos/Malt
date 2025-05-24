package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/taluos/Malt/pkg/errors"
	"github.com/taluos/Malt/pkg/errors/code"
	"github.com/taluos/Malt/server/rest/rest-fiber/internal"
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
func (a AutoStrategy) AuthFunc() fiber.Handler {
	return func(c fiber.Ctx) error {
		operator := AuthOperator{}
		authHeader := strings.SplitN(c.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			internal.WriteResponse(
				c,
				errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Drop()

			return errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong.")
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case "Bearer":
			operator.SetStrategy(a.jwt)
			// a.JWT.MiddlewareFunc()(c)
		default:
			internal.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil)
			c.Drop()

			return errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header.")
		}

		operator.AuthFunc()(c)
		c.Next()

		return nil
	}
}
