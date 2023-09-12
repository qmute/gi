package gi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// LoadFunc 加载用户并将其设置到 gin context 中
type LoadFunc func(c *gin.Context, uid int) error

// MidFake 用于开发环境，模拟任何用户。 生产环境只能在本机调用
// 只需在 request header里加上 Fake-Id ， 值为想要伪装的合法用户ID（系统后台、店铺后台、前端的用户均可）
// 在不同模块里需要传入具体的LoadFunc，完成用户加载
// 为消除安全隐患，在生产环境里只允许 127.0.0.1 使用此功能， 其它环境不限制
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
