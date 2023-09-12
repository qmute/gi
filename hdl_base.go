package gi

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
)

type BaseHdl struct {
}

func GetContext(c *gin.Context) context.Context {
	return c
}

func (p *BaseHdl) GetContext(c *gin.Context) context.Context {
	return GetContext(c)
}

func (p *BaseHdl) Copy(toValue interface{}, fromValue interface{}) error {
	err := copier.Copy(toValue, fromValue)
	if err != nil {
		return WrapInternalCusError(err, "内部错误")
	}

	return nil
}

func (p *BaseHdl) ParseIntQuery(c *gin.Context, key string) (int, bool) {
	s := c.Query(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.WithField("key", key).WithField("value", s).Errorln("bad int query")
		c.String(http.StatusBadRequest, "参数错误")
		c.Abort()
		return 0, false
	}

	return i, true
}

func (p *BaseHdl) HandleError(c *gin.Context, err error, lgs ...*log.Entry) bool {
	return HandleError(c, err, lgs...)
}

func (p *BaseHdl) ParseIntParam(c *gin.Context, key string) (int, bool) {
	s := c.Param(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.WithField("key", key).WithField("value", s).Errorln("bad int param")
		c.String(http.StatusBadRequest, "参数错误")
		c.Abort()
		return 0, false
	}

	return i, true
}

func (p *BaseHdl) Binding(c *gin.Context, obj interface{}, b ...binding.Binding) bool {
	var err error
	if len(b) == 0 {
		err = c.Bind(obj)
	} else {
		err = c.MustBindWith(obj, b[0])
	}

	if err == nil {
		return true
	}

	log.WithError(err).
		WithField("requestId", GetRequestId(c)).
		Errorln("bind error")

	var msg string
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		ret := Translate(errs)
		var arr []string
		for _, v := range ret {
			arr = append(arr, lowerFirst(v))
		}

		msg = strings.Join(arr, "\n")
	} else {
		msg = err.Error()
	}

	c.String(http.StatusBadRequest, msg)
	c.Abort()
	return false

}

// lowerFirst 首字母转小写
func lowerFirst(str string) string {
	if len(strings.TrimSpace(str)) == 0 {
		return str
	}

	r := []rune(str)
	r[0] = unicode.ToLower(r[0])

	// for i, v := range str {
	// 	return string(unicode.ToLower(v)) + str[i+1:]
	// }
	return string(r)
}

func (p *BaseHdl) Valid(c *gin.Context, v Validator) bool {
	if err := v.Valid(); err != nil {
		log.WithError(err).Errorln("valid form error")

		e, ok := IsCusError(err)
		if ok {
			c.String(http.StatusBadRequest, e.Msg())
		} else {
			c.String(http.StatusBadRequest, err.Error())
		}

		c.Abort()
		return false
	}

	return true
}

type Validator interface {
	Valid() error
}
