package config

import (
	"fe_exam/service/validator"
)

func init() {
	initConfig()
	initLogger()
	validator.InitValidator(Config.AppLanguage)
}
