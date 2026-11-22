package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var Config struct {
	AppProd      bool
	AppMode      string
	AppSecret    string
	AppLanguage  string
	MysqlHost    string
	MysqlPort    string
	MysqlName    string
	MysqlUser    string
	MysqlPass    string
	AllowOrigins string
	AllowHeaders string
	LogLevel     string
	AppUploadDir string
}

func envOr(env string, or string) string {
	rt := strings.TrimSpace(os.Getenv(env))
	if rt != "" {
		return rt
	}
	return or
}

func initConfig() {
	Config.AppProd = os.Getenv("APP_PROD") != ""
	if Config.AppProd {
		Config.AppMode = "release"
	} else {
		Config.AppMode = "debug"
	}
	Config.AppSecret = envOr("APP_SECRET", "gin-example:secret")
	Config.AppLanguage = envOr("APP_LANGUAGE", "en")
	Config.MysqlHost = envOr("APP_MYSQL_HOST", "127.0.0.1")
	Config.MysqlPort = envOr("APP_MYSQL_PORT", "3306")
	Config.MysqlName = envOr("APP_MYSQL_NAME", "be_exam")
	Config.MysqlUser = envOr("APP_MYSQL_USER", "root")
	Config.MysqlPass = envOr("APP_MYSQL_PASS", "123456")
	Config.AllowOrigins = envOr("APP_ALLOW_ORIGINS", "http://127.0.0.1:5173")
	Config.AllowHeaders = envOr("APP_ALLOW_HEADERS", "Origin|Content-Length|Content-Type|Authorization")
	Config.LogLevel = envOr("APP_LOG_LEVEL", "info")
	Config.AppUploadDir = envOr("APP_UPLOAD_DIR", "./update")
}
