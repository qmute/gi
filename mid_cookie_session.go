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

type Kv struct {
	Key   interface{}
	Value interface{}
}

func SessionSet(c *gin.Context, key, val interface{}) error {
	se := sessions.Default(c)
	se.Set(key, val)
	return se.Save()
}

func SessionGet(c *gin.Context, key interface{}) interface{} {
	se := sessions.Default(c)
	v := se.Get(key)
	return v
}

func SessionDelete(c *gin.Context, key interface{}) error {
	se := sessions.Default(c)
	se.Delete(key)
	return se.Save()
}

func SessionBatchSet(c *gin.Context, kvs ...Kv) error {
	if len(kvs) == 0 {
		return nil
	}

	se := sessions.Default(c)
	for _, kv := range kvs {
		se.Set(kv.Key, kv.Value)
	}

	return se.Save()
}

func SessionBatchDelete(c *gin.Context, keys ...interface{}) error {
	if len(keys) == 0 {
		return nil
	}

	se := sessions.Default(c)
	for _, k := range keys {
		se.Delete(k)
	}

	return se.Save()
}
