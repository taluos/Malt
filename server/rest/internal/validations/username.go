package validations

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterUsername(trans uTranslator.Translator) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("username", validateUsername)
		_ = v.RegisterTranslation("username", trans,
			func(ut uTranslator.Translator) error {
				return ut.Add("username", "{0} 格式错误", true)
			},
			func(ut uTranslator.Translator, fe validator.FieldError) string {
				t, _ := ut.T("username", fe.Field())
				return t
			})
	}
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// 用户名只能包含字母、数字、下划线，长度3-16位
	ok, err := regexp.MatchString(`^[a-zA-Z0-9_]{3,16}$`, username)
	if err != nil {
		return false
	}
	return ok
}
