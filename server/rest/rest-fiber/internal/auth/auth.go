// Basic Auth middleware
package auth

import (
	fiber "github.com/gofiber/fiber/v3"
)

// AuthStrategy defines the set of methods used to do resource authentication.
type AuthStrategy interface {
	AuthFunc() fiber.Handler
}

// AuthOperator used to switch between different authentication strategy.
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy used to set to another authentication strategy.
func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

// AuthFunc execute resource authentication.
func (operator *AuthOperator) AuthFunc() fiber.Handler {
	return operator.strategy.AuthFunc()
}
