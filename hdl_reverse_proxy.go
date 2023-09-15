package gi

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

// NewReverseProxyHandler 创建反向代理handler，当代理请求不能正常执行时，返回404，隐藏后端oss详细信息
func NewReverseProxyHandler(target *url.URL) gin.HandlerFunc {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			uri := r.Out.URL
			lastPath, _ := lo.Last(strings.Split(uri.Path, "/"))
			// path最后一部分没有"."，说明不是文件，是文件夹，则为其添加index.html（未作详细文件类型匹配）
			if !strings.Contains(lastPath, ".") {
				uri.Path, _ = url.JoinPath(uri.Path, "index.html")
				uri.RawPath, _ = url.JoinPath(uri.RawPath, "index.html")
			}
		},
		ModifyResponse: func(r *http.Response) error {
			// 不正常
			if r.StatusCode >= 400 {
				b, _ := io.ReadAll(r.Body)
				logrus.Warningln("proxy fail", string(b))
				return errors.Errorf("not found")
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
		},
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorln("proxy err", err)
				c.AbortWithStatus(http.StatusServiceUnavailable) // 当作临时不可用
			}
		}()
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
