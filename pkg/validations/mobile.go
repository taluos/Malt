package validations

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterMobile(trans uTranslator.Translator) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validateMobile)
		_ = v.RegisterTranslation("mobile", trans,
			func(ut uTranslator.Translator) error {
				return ut.Add("mobile", "{0} 格式错误", true)
			},
			func(ut uTranslator.Translator, fe validator.FieldError) string {
				t, _ := ut.T("mobile", fe.Field())
				return t
			})
	}
}

func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()

	ok, err := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if err != nil {
		return false
	}
	return ok
}
