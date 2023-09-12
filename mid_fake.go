package gi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LoadFunc 加载用户并将其设置到 gin context 中
type LoadFunc func(c *gin.Context, uid int) error

// MidFake 用于开发环境，模拟任何用户。 生产环境只能在本机调用
func MidFake(f LoadFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.GetHeader("Fake-Id")

		// 常规请求
		if s == "" {
			return
		}

		id, err := strconv.Atoi(s)
		if err != nil {
			log.WithError(err).Warnln("bad fake id", s)
			return
		}

		if gin.Mode() == gin.ReleaseMode && c.ClientIP() != "127.0.0.1" {
			return
		}
		if err := f(c, id); err != nil {
			log.WithError(err).Errorln("load fake user err")
			return
		}
	}
}
