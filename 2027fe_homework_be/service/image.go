package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fe_exam/config"
)

type ImageSrv struct {
}

const (
	maxImageSize = 5 << 20 // 5MB

	// ImageURLPath 静态访问前缀；router 挂载静态目录时复用这个常量，避免两边写死不一致
	ImageURLPath   = "/update"
	imageURLPrefix = ImageURLPath + "/"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// Save 保存上传图片到 Config.AppUploadDir，返回相对路径（如 /image/20260101_x.png）
func (i *ImageSrv) Save(r io.Reader, filename string, size int64) (ImageUploadResponse, error) {
	if size <= 0 {
		return ImageUploadResponse{}, ErrParam.WithMsg("图片不能为空")
	}
	if size > maxImageSize {
		return ImageUploadResponse{}, ErrParam.WithMsg("图片过大，最大 5MB")
	}

	buf, err := io.ReadAll(io.LimitReader(r, maxImageSize+1))
	if err != nil {
		return ImageUploadResponse{}, ErrInternal.With(err)
	}
	if len(buf) == 0 {
		return ImageUploadResponse{}, ErrParam.WithMsg("图片不能为空")
	}
	if len(buf) > maxImageSize {
		return ImageUploadResponse{}, ErrParam.WithMsg("图片过大，最大 5MB")
	}

	contentType := http.DetectContentType(buf)
	ext, ok := allowedImageTypes[contentType]
	if !ok {
		return ImageUploadResponse{}, ErrParam.WithMsg("仅支持 jpg/png/webp/gif 格式")
	}

	dir := config.Config.AppUploadDir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ImageUploadResponse{}, ErrInternal.With(err)
	}

	name := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), randomSuffix(), ext)
	if err := os.WriteFile(filepath.Join(dir, name), buf, 0o644); err != nil {
		return ImageUploadResponse{}, ErrInternal.With(err)
	}

	return ImageUploadResponse{URL: imageURLPrefix + name}, nil
}

func randomSuffix() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
