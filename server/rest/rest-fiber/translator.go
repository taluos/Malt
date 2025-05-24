package fiber

import (
	"github.com/taluos/Malt/pkg/errors"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func initTrans(locale string) (ut.Translator, error) {
	// 修改为包内的变量名
	en := en.New()
	zh := zh.New()

	uni := ut.New(en, zh)

	trans, ok := uni.GetTranslator(locale)
	if !ok {
		return nil, errors.Errorf("uni.GetTranslator(%s) failed", locale)
	}

	// 注册翻译器
	v := validator.New()
	switch locale {
	case "en":
		err := en_translations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			return nil, errors.Errorf("en_translations.RegisterDefaultTranslations failed: %s", err)
		}
	case "zh":
		err := zh_translations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			return nil, errors.Errorf("zh_translations.RegisterDefaultTranslations failed: %s", err)
		}
	default:
		err := en_translations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			return nil, errors.Errorf("translations.RegisterDefaultTranslations failed: %s", err)
		}
	}
	return trans, nil
}
