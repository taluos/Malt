package validations

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterEmail(trans uTranslator.Translator) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("email", validateEmail)
		_ = v.RegisterTranslation("email", trans,
			func(ut uTranslator.Translator) error {
				return ut.Add("email", "{0} 格式错误", true)
			},
			func(ut uTranslator.Translator, fe validator.FieldError) string {
				t, _ := ut.T("email", fe.Field())
				return t
			})
	}
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	// 简单的邮箱正则表达式
	ok, err := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil {
		return false
	}
	return ok
}
