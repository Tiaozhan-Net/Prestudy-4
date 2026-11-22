package router

import (
	"net/http"
	"os"

	"fe_exam/config"
	"fe_exam/controller"
	"fe_exam/logger"
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

func NewServer() *http.Server {
	r := gin.Default()
	config.SetCORS(r)
	config.InitSession(r)

	// 上传图片目录：Config.AppUploadDir（APP_UPLOAD_DIR，默认 ./update），供 /update/xxx.png 静态访问
	if err := os.MkdirAll(config.Config.AppUploadDir, 0o755); err != nil {
		logger.Warnf("create upload dir %q failed: %v", config.Config.AppUploadDir, err)
	}
	r.Static(service.ImageURLPath, config.Config.AppUploadDir)
	InitRouter(r)
	s := &http.Server{
		Addr:    "0.0.0.0:8088",
		Handler: r,
	}
	return s

}

var ctr = controller.New()
