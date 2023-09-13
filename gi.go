package gi

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
)

type GinOption func(*gin.Engine)

// New 创建 gin.Engine, 可指定多个Option
func New(opt ...GinOption) *gin.Engine {
	binding.Validator = new(defaultValidator)

	if err := initTrans(ZH); err != nil {
		log.WithError(err).Errorln("init trans failed")
	}

	if !gin.IsDebugging() {
		gin.DisableConsoleColor()
	}
	r := gin.New()
	for _, v := range opt {
		v(r)
	}
	return r
}

// WithStatic 服务静态文件
// 默认url前缀为 /，本地文件路径为 ./public，自动索引
// 如果需要自定义 ，可使用 gi.Static middleware
func WithStatic() GinOption {
	return with(Static(StaticWithIndex(true)))
}

// WithPprof 启用pprof
func WithPprof() GinOption {
	return func(router *gin.Engine) {
		pprof.Register(router)
	}
}

// With 使用任意middleware
func with(fn gin.HandlerFunc) GinOption {
	return func(r *gin.Engine) {
		r.Use(fn)
	}
}
