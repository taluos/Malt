package validations

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterPassword(trans uTranslator.Translator) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("password", validatePassword)
		_ = v.RegisterTranslation("password", trans,
			func(ut uTranslator.Translator) error {
				return ut.Add("password", "{0} 格式错误，密码需8-20位且包含字母和数字，且仅允许使用部分特殊字符(!@#$%^&*_-)", true)
			},
			func(ut uTranslator.Translator, fe validator.FieldError) string {
				t, _ := ut.T("password", fe.Field())
				return t
			})
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 密码8-20位，必须包含字母和数字，只允许字母、数字和!@#$%^&*_-这几个特殊字符
	ok, err := regexp.MatchString(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*_-]{8,20}$`, password)
	if err != nil {
		return false
	}
	return ok
}
