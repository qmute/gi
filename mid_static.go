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

// StaticWithIndex 是否自动索引，默认为 false
func StaticWithIndex(index bool) StaticOption {
	return func(cfg *staticConfig) {
		cfg.index = index
	}
}

// Static 静态文件服务
func Static(opt ...StaticOption) gin.HandlerFunc {
	cfg := &staticConfig{
		urlPrefix: "/",
		fileRoot:  "./public",
		index:     false,
	}
	for _, v := range opt {
		v(cfg)
	}
	return static.Serve(cfg.urlPrefix, static.LocalFile(cfg.fileRoot, cfg.index))
}
