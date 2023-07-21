package genericmiddleware

import (
	"net/http"
	"time"

	gindump "github.com/tpkeeper/gin-dump"

	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/pkg/sets"
)

type SkipperFunc func(*gin.Context) bool

func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		return sets.NewString(prefixes...).IsPrefixOf(path)
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		return !sets.NewString(prefixes...).IsPrefixOf(path)
	}
}

// SkipHandler 跳过指定的中间件.
func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

// Middlewares store registered middlewares.
var (
	MiddlewareList  = defaultMiddlewareList()
	MiddlewareNames = defaultMiddlewareListNames()
)

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")

	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}

func defaultMiddlewareList() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"context":   Context(),
		"requestid": RequestID(),
		"recovery":  gin.Recovery(),
		"secure":    Secure,
		"options":   Options,
		"nocache":   NoCache,
		"cors":      Cors(),
		"logger":    Logger(),
		"dump":      gindump.Dump(),
	}
}

func defaultMiddlewareListNames() []string {
	names := make([]string, 0, len(defaultMiddlewareList()))
	for name := range defaultMiddlewareList() {
		names = append(names, name)
	}
	return names
}
