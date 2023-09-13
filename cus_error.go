package gi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/quexer/utee"
	log "github.com/sirupsen/logrus"
	micErr "go-micro.dev/v4/errors"
	"gorm.io/gorm"
)

const (
	ErrCodeOk           ErrCode = 0
	ErrCodeNotModified  ErrCode = 304
	ErrCodeBadReq       ErrCode = 400
	ErrCodeUnauthorized ErrCode = 401
	ErrCodeForbidden    ErrCode = 403
	ErrCodeNotFound     ErrCode = 404
	ErrCodeInternalErr  ErrCode = 500
	ErrCodePanicErr     ErrCode = 590 // internal error, but panic error
)

type ErrCode int

func (p ErrCode) IsOk() bool {
	return p == ErrCodeOk
}

func (p ErrCode) Value() int {
	return int(p)
}

type CusError struct {
	code    ErrCode
	msg     string
	context utee.J
	err     error
}

func (p *CusError) Error() string {
	if p.err != nil {
		return fmt.Sprintf("%s:%+v:%+v", p.msg, p.context, p.err)
	}
	return p.msg
}

func (e *CusError) Unwrap() error {
	return e.err
}

func (p *CusError) Msg() string {
	return p.msg
}

func (p *CusError) Code() ErrCode {
	return p.code
}

func (p *CusError) Context() utee.J {
	return p.context
}

func NewCusError(code ErrCode, msg string, contexts ...utee.J) error {
	err := &CusError{
		code: code,
		msg:  msg,
	}
	if len(contexts) > 0 {
		err.context = contexts[0]
	}
	return errors.WithStack(err)
}

func WrapCusErr(code ErrCode, e error, msg string, contexts ...utee.J) error {
	err := &CusError{
		code: code,
		msg:  msg,
		err:  e,
	}
	if len(contexts) > 0 {
		err.context = contexts[0]
	}
	return err
}

func WrapBadRequestCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodeBadReq, err, msg, contexts...)
}

func WrapNotFoundCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodeNotFound, err, msg, contexts...)
}

func WrapForbiddenCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodeForbidden, err, msg, contexts...)
}

func WrapUnauthorizedCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodeUnauthorized, err, msg, contexts...)
}

func WrapInternalCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodeInternalErr, err, msg, contexts...)
}

func WrapPanicCusError(err error, msg string, contexts ...utee.J) error {
	return WrapCusErr(ErrCodePanicErr, err, msg, contexts...)
}

func IsCusError(err error) (*CusError, bool) {
	if err == nil {
		return nil, false
	}

	e, ok := err.(*CusError)
	if !ok {
		return nil, ok
	}

	return e, true
}

func GetErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := IsCusError(err); ok {
		return e.Msg()
	}

	return err.Error()
}

func HandleError(c *gin.Context, err error, lgs ...*log.Entry) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.String(http.StatusBadRequest, "没有找到记录")
		c.Abort()
		return true
	}

	// 测试模式下保持安静
	if gin.Mode() != gin.TestMode {
		fmt.Printf("%+v", err) // 打印到标准输出，方便查错
	}

	mErr := micErr.Parse(err.Error())
	if mErr.GetCode() > 0 {
		err = WrapCusErr(ErrCode(mErr.GetCode()), err, mErr.GetDetail())
	}

	ce, ok := IsCusError(err)
	if !ok {
		// 如果是未包装错误， 现在包装， 下面统一处理
		ce = WrapInternalCusError(err, "服务错误，请稍后重试").(*CusError)
	}

	clientIP := c.ClientIP()
	requestId := GetRequestId(c)
	lg := errEntry(ce, lgs...).WithField("ip", clientIP).WithField("requestId", requestId)

	for k, v := range ce.Context() {
		lg = lg.WithField(k, v)
	}

	outMsg := func() string {
		if ce.Code() >= 500 {
			return "panic " + ce.Msg()
		} else {
			return ce.Msg()
		}
	}

	lg.Errorln(outMsg())

	if ce.Code() > 500 {
		c.String(http.StatusInternalServerError, ce.Msg())
		c.Abort()
		return true
	}

	c.String(int(ce.Code()), ce.Msg())
	c.Abort()
	return true
}

func errEntry(err error, vs ...*log.Entry) *log.Entry {
	msg := fmt.Sprintf("%+v", err)
	if len(vs) > 0 {
		return vs[0].WithField("err", msg)
	}
	return log.WithField("err", msg)
}
