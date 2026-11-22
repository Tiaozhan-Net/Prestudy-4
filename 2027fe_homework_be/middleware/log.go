package middleware

import (
	"bytes"
	"fe_exam/config"
	"fe_exam/logger"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// maxCaptureBody 超过这个大小的请求体不再读取，避免图片/大 JSON 进内存
	maxCaptureBody = 64 << 10
	// maxLogBody 日志里 body 最多保留这么多字节，其余截断
	maxLogBody = 4 << 10
)

// isTextContentType 只有文本类请求体才值得落日志
func isTextContentType(contentType string) bool {
	if contentType == "" {
		return true
	}
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	switch {
	case strings.HasPrefix(ct, "text/"),
		strings.HasPrefix(ct, "application/json"),
		strings.HasPrefix(ct, "application/xml"),
		ct == "application/x-www-form-urlencoded",
		ct == "application/graphql":
		return true
	default:
		return false
	}
}

// bodyPreview 生成可安全写入日志的 body 摘要：二进制只记元信息，文本截断并修正非法 UTF-8
func bodyPreview(contentType string, size int64, body []byte) string {
	if !isTextContentType(contentType) {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
		if mediaType == "" {
			mediaType = "unknown"
		}
		return fmt.Sprintf("<%s, %d bytes, binary skipped>", mediaType, size)
	}
	if len(body) == 0 {
		return ""
	}

	truncated := len(body) > maxLogBody
	preview := body
	if truncated {
		preview = body[:maxLogBody]
	}
	s := strings.ToValidUTF8(string(preview), ".")
	if truncated {
		s += "...(truncated)"
	}
	return s
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		w := &logger.ResponseBodyWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 只缓存文本且不超过 64KB 的请求体；二进制（上传的图片等）不读，避免污染日志并省内存
		var requestBody []byte
		requestContentType := c.Request.Header.Get("Content-Type")
		requestContentLength := c.Request.ContentLength
		if c.Request.Body != nil && isTextContentType(requestContentType) &&
			(requestContentLength < 0 || requestContentLength <= maxCaptureBody) {
			bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, maxCaptureBody+1))
			if len(bodyBytes) > 0 {
				// 把读过的部分拼回去，保证 handler 仍能读到完整 body
				c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), c.Request.Body))
			}
			if err == nil {
				requestBody = bodyBytes
			}
		}

		// Create a copy of the context for logging
		// logContext := *c
		logContext := struct {
			Writer   gin.ResponseWriter
			Request  *http.Request
			ClientIP func() string
		}{
			Writer:   c.Writer,
			Request:  c.Request,
			ClientIP: c.ClientIP,
		}

		c.Next()
		go func() {
			select {
			case <-config.SkipSignalChan:
				return
			default:
				goto log
			}
		log:
			status := logContext.Writer.Status()
			path := logContext.Request.URL.Path
			query := logContext.Request.URL.RawQuery
			cost := time.Since(start)
			method := logContext.Request.Method
			clientIP := logContext.ClientIP()
			userAgent := logContext.Request.UserAgent()

			if logger.GinLogger.Level == logrus.DebugLevel {
				if logContext.Writer != nil {
					responseHeaders := logContext.Writer.Header()
					responseBody := w.Body.Bytes()
					responseContentType := responseHeaders.Get("Content-Type")
					requestHeaders, _ := httputil.DumpRequest(logContext.Request, false)

					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":           method,
						"\nurl":              path,
						"\nquery":            query,
						"\nclient_ip":        clientIP,
						"\nuser_agent":       userAgent,
						"\nstatus":           status,
						"\nduration":         cost,
						"\nrequest_headers":  string(requestHeaders),
						"\nrequest_body":     bodyPreview(requestContentType, requestContentLength, requestBody),
						"\nresponse_headers": responseHeaders,
						"\nresponse_body":    bodyPreview(responseContentType, int64(len(responseBody)), responseBody),
					}).Debug()
				}
			} else {
				switch {
				case status >= http.StatusInternalServerError:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip":  clientIP,
						"\nuser_agent": userAgent,
						"\nStatus":     status,
						"\nduration":   cost}).Error()
				case status >= http.StatusBadRequest:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip":  clientIP,
						"\nuser_agent": userAgent,
						"\nstatus":     status,
						"\nduration":   cost}).Warn()
				default:
					logger.GinLogger.WithFields(logrus.Fields{
						"\nmethod":     method,
						"\nurl":        path,
						"\nquery":      query,
						"\nclient_ip:": clientIP,
						"\nuser_agent": userAgent,
						"\nstatus":     status,
						"\nduration":   cost}).Info()
				}
			}
		}()
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.GinLogger.Error("broken pipe: ", err, ". Request: ", string(httpRequest))
					c.Abort()
					return
				}

				//deeper stack
				pc, file, line, _ := runtime.Caller(3)
				func_name := runtime.FuncForPC(pc).Name()

				if stack {
					logger.GinLogger.Errorf("panic recovered: %v. Request: %s. File: %s, Line: %d, Function: %s. Stack: %s", err, string(httpRequest), file, line, func_name, string(debug.Stack()))
				} else {
					logger.GinLogger.Errorf("panic recovered: %v. Request: %s. File: %s, Line: %d, Function: %s", err, string(httpRequest), file, line, func_name)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
