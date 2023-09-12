package gi

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type GinOption func(*gin.Engine)

// New 创建 gin.Engine, 可指定多个Option
func New(opt ...GinOption) *gin.Engine {
	if !gin.IsDebugging() {
		gin.DisableConsoleColor()
	}
	r := gin.New()
	for _, v := range opt {
		v(r)
	}
	return r
}

// WithCookieSession CookieSession middleware
// name, cookie name.
// salt, cookie store secret
func WithCookieSession(name, salt string) GinOption {
	store := cookie.NewStore([]byte(salt))
	return With(sessions.Sessions(name, store))
}

// WithStatic 服务静态文件，fileRoot 本地文件路径，默认为 ./public
func WithStatic(fileRoot ...string) GinOption {
	root := "./public"
	if len(fileRoot) > 0 {
		root = fileRoot[0]
	}
	return With(static.ServeRoot("/", root))
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
