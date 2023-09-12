package gi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// MidCookieSession CookieSession middleware
// name, cookie name.
// salt, cookie store secret
// opt, optional session option
func MidCookieSession(name, salt string, opt ...sessions.Options) gin.HandlerFunc {
	store := cookie.NewStore([]byte(salt))
	if len(opt) > 0 {
		store.Options(opt[0])
	}
	return sessions.Sessions(name, store)
}
