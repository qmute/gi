package gi

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
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

// WithStatic 服务静态文件，fileRoot 本地文件路径，默认文件为 ./public
// 对于目录，自动索引
func WithStatic(fileRoot ...string) GinOption {
	root := "./public"
	if len(fileRoot) > 0 {
		root = fileRoot[0]
	}
	return With(static.Serve("/", static.LocalFile(root, true)))
}

// WithPprof 启用pprof
func WithPprof() GinOption {
	return func(router *gin.Engine) {
		pprof.Register(router)
	}
}

// With 使用任意middleware
func With(fn gin.HandlerFunc) GinOption {
	return func(r *gin.Engine) {
		r.Use(fn)
	}
}
