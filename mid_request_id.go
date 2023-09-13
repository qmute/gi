package gi

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// MidRequestId request id middleware
var MidRequestId = requestid.New

// GetRequestId get request id from context
func GetRequestId(c *gin.Context) string {
	return requestid.Get(c)
}
