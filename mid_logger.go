package gi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

// LogOpt Gin日志配置项
type LogOpt func(opt *logConfig)

// FieldGetter 从 gin.Context 中获取内容的函数
type FieldGetter func(r *http.Request) any

type logConfig struct {
	Threshold time.Duration

	FieldGetter map[string]FieldGetter

	IgnoreStaticMedia bool // 是否忽略静态资源
}

// LogWithIgnoreStaticMedia 是否忽略静态资源，默认为true，忽略
// 静态资源由常见的后缀判断：.css, .js, .html, .png, .gif, .jpg, .jpeg, .ico
func LogWithIgnoreStaticMedia(ignore bool) LogOpt {
	return func(opt *logConfig) {
		opt.IgnoreStaticMedia = ignore
	}
}

// LogWithThreshold 设置日志打印时限，超过时限才打印日志. 默认为0，即不限制
func LogWithThreshold(threshold time.Duration) LogOpt {
	return func(opt *logConfig) {
		opt.Threshold = threshold
	}
}

// LogWithField 为日志提供自定义字段
// name: 字段名
// getter: 字段值获取函数
func LogWithField(name string, getter FieldGetter) LogOpt {
	return func(opt *logConfig) {
		opt.FieldGetter[name] = getter
	}
}

// MidLogger 打印日志
func MidLogger(opt ...LogOpt) gin.HandlerFunc {
	config := &logConfig{
		Threshold:         0,
		FieldGetter:       map[string]FieldGetter{},
		IgnoreStaticMedia: true,
	}
	for _, v := range opt {
		v(config)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if config.IgnoreStaticMedia && ignorePath(path) {
			return
		}

		buf, err := io.ReadAll(c.Request.Body)
		utee.Chk(err)
		bodyCopyReader := io.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		method := c.Request.Method
		if method == "HEAD" {
			return
		}

		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		entry := log.WithField("mod", "gin").
			WithField("latency", latency.String()).
			WithField("ip", c.ClientIP()).
			WithField("method", method).
			WithField("path", path).
			WithField("lat", fmt.Sprintf("%.2f", float64(latency.Nanoseconds())/1e6)). // 单位为毫秒
			WithField("ua", c.Request.Header.Get("user-agent")).
			WithField("requestId", requestid.Get(c))

		for k, v := range config.FieldGetter {
			entry = entry.WithField(k, v(c.Request))
		}

		if comment != "" {
			entry = entry.WithField("err", comment)
		}

		if statusCode != http.StatusOK && statusCode != http.StatusNotModified {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(bodyCopyReader); err != nil {
				errEntry(err, entry).Errorln("mid error")
			} else {
				entry = entry.WithField("body", buf.String()).WithField("query", c.Request.URL.RawQuery)
			}

			return
		}

		// 没传时限或时限已超，则打印日志
		if config.Threshold == 0 || latency > config.Threshold {
			entry.Infoln(statusCode)
		}
	}
}

// 要忽略的静态资源后缀
var ignoreFileExtension = []string{".css", ".js", ".html", ".png", ".gif", ".jpg", ".jpeg", ".ico"}

func ignorePath(path string) bool {
	fn := func(x string) bool {
		return strings.HasSuffix(path, x)
	}
	return lo.SomeBy(ignoreFileExtension, fn)
}
