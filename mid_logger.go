package gi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/quexer/utee"
	log "github.com/sirupsen/logrus"

	"baibao/meishi/pkg/ut"
)

// MidLogger 出错时打印日志
func MidLogger(threshold ...time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		buf, err := ioutil.ReadAll(c.Request.Body)
		utee.Chk(err)
		bodyCopyReader := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if method == "HEAD" {
			return
		}

		platform := ut.GetXPlatform(c.Request)

		ua := c.Request.Header.Get("user-agent")

		h5Tid := 0
		if strings.HasPrefix(path, "/shop/") {
			str := c.Request.URL.Query().Get("tid")
			if str != "" {
				h5Tid, _ = strconv.Atoi(str)
			}
		}

		reqId := requestid.Get(c)

		entry := log.WithField("mod", "gin").
			WithField("platform", platform).
			WithField("latency", latency.String()).
			WithField("ip", clientIP).WithField("method", method).
			WithField("path", path).
			WithField("lat", fmt.Sprintf("%.2f", float64(latency.Nanoseconds())/1e6)). // 单位为毫秒
			WithField("ua", ua).WithField("requestId", reqId)

		if h5Tid != 0 {
			entry = entry.WithField("h5tid", h5Tid)
		}

		if comment != "" {
			entry = entry.WithField("err", comment)
		}

		if statusCode != http.StatusOK && statusCode != http.StatusNotModified {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(bodyCopyReader); err != nil {
				ut.ErrEntry(err, entry).Errorln("mid error")
			} else {
				entry = entry.WithField("body", buf.String()).WithField("query", c.Request.URL.RawQuery)
			}

			return
		}

		// 没传时限或时限已超，打印日志
		if len(threshold) == 0 || latency > threshold[0] {
			entry.Infoln(statusCode)
		}
	}
}
