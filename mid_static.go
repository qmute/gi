package gi

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type staticConfig struct {
	urlPrefix string
	fileRoot  string
	index     bool
}

// StaticOption 初始化静态文件服务的option
type StaticOption func(*staticConfig)

// StaticWithUrlPrefix 静态文件服务的url前缀，默认为 /
func StaticWithUrlPrefix(urlPrefix string) StaticOption {
	return func(cfg *staticConfig) {
		cfg.urlPrefix = urlPrefix
	}
}

// StaticWithFileRoot 静态文件服务的本地文件路径，默认为 ./public
func StaticWithFileRoot(fileRoot string) StaticOption {
	return func(cfg *staticConfig) {
		cfg.fileRoot = fileRoot
	}
}

// StaticWithoutIndex 禁用自动索引
func StaticWithoutIndex() StaticOption {
	return func(cfg *staticConfig) {
		cfg.index = false
	}
}

// Static 静态文件服务
func Static(opt ...StaticOption) gin.HandlerFunc {
	cfg := &staticConfig{
		urlPrefix: "/",
		fileRoot:  "./public",
		index:     true,
	}
	for _, v := range opt {
		v(cfg)
	}
	return static.Serve(cfg.urlPrefix, static.LocalFile(cfg.fileRoot, cfg.index))
}
