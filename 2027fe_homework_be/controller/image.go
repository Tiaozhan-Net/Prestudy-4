package controller

import (
	"fe_exam/service"

	"github.com/gin-gonic/gin"
)

type Image struct {
}

// Upload POST /images
func (i *Image) Upload(c *gin.Context) {
	if _, logged := currentUID(c); !logged {
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		fail(c, service.ErrParam.WithMsg("缺少 image 文件字段"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		fail(c, service.ErrInternal.With(err))
		return
	}
	defer file.Close()

	resp, err := srv.ImageSrv.Save(file, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, resp)
}
