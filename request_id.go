package gi

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// RequestId request id middleware
var RequestId = requestid.New

// GetRequestId get request id from context
func GetRequestId(c *gin.Context) string {
	return requestid.Get(c)
}
