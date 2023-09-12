package gi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"
	log "github.com/sirupsen/logrus"
)

// MidLogger 出错时打印日志
func MidLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bodyLogWriter

	buf, err := io.ReadAll(c.Request.Body)
	utee.Chk(err)
	bodyCopyReader := io.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

	c.Next()

	latency := time.Since(start)

	clientIP := c.ClientIP()
	method := c.Request.Method
	header := c.Request.Header
	statusCode := c.Writer.Status()
	comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
	if method == "HEAD" {
		return
	}

	entry := log.WithField("mod", "gin").
		WithField("latency", latency.String()).
		WithField("ip", clientIP).
		WithField("method", method).
		WithField("path", path).
		WithField("lat", fmt.Sprintf("%.2f", float64(latency.Nanoseconds())/1e6))
	if comment != "" {
		entry = entry.WithField("err", comment)
	}

	if statusCode != http.StatusOK && statusCode != http.StatusNotModified {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(bodyCopyReader); err != nil {
			entry.WithError(err).Errorln("mid error")
		} else {
			if bodyLogWriter.body != nil {
				entry = entry.WithField("response", bodyLogWriter.body.String())
			}
			entry = entry.WithField("header", header).
				WithField("body", buf.String()).
				WithField("query", c.Request.URL.RawQuery)
		}
	}
	entry.Infoln(statusCode)
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
