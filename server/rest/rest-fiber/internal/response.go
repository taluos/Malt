package internal

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/taluos/Malt/pkg/errors"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"msg"`

	Detail string `json:"detail"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c fiber.Ctx, err error, data interface{}) {
	if err != nil {
		errStr := fmt.Sprintf("%#+v", err)
		coder := errors.ParseCoder(err)
		c.Status(coder.HTTPStatus())
		c.JSON(ErrResponse{
			Code:      coder.Code(),
			Message:   coder.String(),
			Detail:    errStr,
			Reference: coder.Reference(),
		})

		return
	}
	c.Status(http.StatusOK)
	c.JSON(data)
}
