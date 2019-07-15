package kuu

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Context
type Context struct {
	*gin.Context
	SignInfo      *SignContext
	PrisDesc      *PrivilegesDesc
	RoutineCaches RoutineCaches
}

// L
func (c *Context) L(defaultValue string, args ...interface{}) string {
	return L(c.Context, defaultValue, args...)
}

// Lang
func (c *Context) Lang(key string, defaultValue string, args interface{}) string {
	return Lang(c.Context, key, defaultValue, args)
}

// DB
func (c *Context) DB() *gorm.DB {
	return DB()
}

// WithTransaction
func (c *Context) WithTransaction(fn func(*gorm.DB) (*gorm.DB, error), with ...*gorm.DB) error {
	return WithTransaction(fn, with...)
}

// STD
func (c *Context) STD(data interface{}, msg ...string) *STDRender {
	return STD(c.Context, data, msg...)
}

// STDErr
func (c *Context) STDErr(msg string, err ...interface{}) *STDRender {
	return STDErr(c.Context, msg, err...)
}

// STDHold
func (c *Context) STDHold(data interface{}, msg ...string) *STDRender {
	return STDHold(c.Context, data, msg...)
}

// STDErrHold
func (c *Context) STDErrHold(msg string, err ...interface{}) *STDRender {
	return STDErrHold(c.Context, msg, err...)
}

// SetValue
func (c *Context) SetRoutineCache(key string, value interface{}) {
	SetRoutineCache(key, value)
}

// DelValue
func (c *Context) DelRoutineCache(key string) {
	DelRoutineCache(key)
}

// GetValue
func (c *Context) GetRoutineCache(key string) interface{} {
	return GetRoutineCache(key)
}

// PRINT
func (c *Context) PRINT(args ...interface{}) {
	PRINT(args...)
}

// DEBUG
func (c *Context) DEBUG(args ...interface{}) {
	DEBUG(args...)
}

// WARN
func (c *Context) WARN(args ...interface{}) {
	WARN(args...)
}

// INFO
func (c *Context) INFO(args ...interface{}) {
	INFO(args...)
}

// ERROR
func (c *Context) ERROR(args ...interface{}) {
	ERROR(args...)
}

// FATAL
func (c *Context) FATAL(args ...interface{}) {
	FATAL(args...)
}

// PANIC
func (c *Context) PANIC(args ...interface{}) {
	PANIC(args...)
}

// IgnoreAuth
func (c *Context) IgnoreAuth(cancel ...bool) *Context {
	c.RoutineCaches.IgnoreAuth(cancel...)
	return c
}

// Scheme
func (c *Context) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.Request.TLS != nil {
		return "https"
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := c.Request.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := c.Request.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}

// Origin
func (c *Context) Origin() string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request.Host)
}
