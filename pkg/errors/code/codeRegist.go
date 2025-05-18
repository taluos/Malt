package code

import (
	"fmt"
	"net/http"

	"Malt/pkg/errors"
	"Malt/pkg/log"

	"github.com/novalagung/gubrak"
)

var IncludeErrCode = []int{200, 400, 401, 403, 404, 500}

type errCode struct {
	// code 错误码
	C int

	// http 状态码
	HTTP int

	// 扩展字段
	Ext string

	// 引用文档
	Ref string
}

func (e errCode) Code() int {
	if e.C == 0 {
		return ErrUnknow
	}
	return e.C
}

func (e errCode) HTTPStatus() int {

	if e.HTTP == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTP
}

func (e errCode) String() string {
	return e.Ext
}

func (e errCode) Reference() string {
	return e.Ref
}

func register(code int, HttpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes(IncludeErrCode, HttpStatus)
	if !found {
		log.Fatal("HTTP code is not available")
	}

	var ref string = ""
	for _, v := range refs {
		ref = fmt.Sprintf(ref, "$$", v)
	}

	var coder = errCode{
		C:    code,
		HTTP: HttpStatus,
		Ext:  message,
		Ref:  ref,
	}
	errors.MustRegister(coder)
}

var _ errors.Coder = (*errCode)(nil)
