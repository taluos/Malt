package httpserver

import (
	"reflect"

	"github.com/taluos/Malt/pkg/log"

	"github.com/taluos/Malt/pkg/errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	uTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_tw "github.com/go-playground/validator/v10/translations/en"
	zh_tw "github.com/go-playground/validator/v10/translations/zh"
)

func initTrans(transtype string) (uTranslator.Translator, error) {
	var (
		err   error
		trans uTranslator.Translator
	)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 获取Json中tag的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := field.Tag.Get("label")
			if name == "" {
				return field.Name
			}
			return name

		})

		zhTrans := zh.New()                               // 中文翻译
		enTrans := en.New()                               // 英文翻译
		uni := uTranslator.New(enTrans, zhTrans, enTrans) // 第一个参数为备用语言环境， 后续参数为应该支持的语言环境

		trans, ok = uni.GetTranslator(transtype)
		if !ok {
			log.Errorf("Error in [uni.GetTranslator] [Translator 初始化错误] .")
			return nil, errors.New("Init Translator error")
		}

		switch transtype {
		case "en":
			err = en_tw.RegisterDefaultTranslations(v, trans)
			if err != nil {
				log.Errorf("Error in [en_translations.RegisterDefaultTranslations] [Translator 初始化错误] .")
				return nil, errors.Wrapf(err, "Init en Translator error")
			}
		case "zh":
			err = zh_tw.RegisterDefaultTranslations(v, trans)
			if err != nil {
				log.Errorf("Error in [zh_translations.RegisterDefaultTranslations] [Translator 初始化错误] .")
				return nil, errors.Wrapf(err, "Init zh Translator error")
			}
		default:
			err = en_tw.RegisterDefaultTranslations(v, trans)
			if err != nil {
				log.Errorf("Error in [en_translations.RegisterDefaultTranslations] [Translator 初始化错误] .")
				return nil, errors.Wrapf(err, "Init Translator error")
			}
		}
	}
	return trans, nil
}
