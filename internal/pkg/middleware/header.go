package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 是一个 Gin 的中间件, 用来禁止客户端缓存 HTTP 请求的返回结果.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	
	// 设置HTTP 的 Expires 响应头告诉客户端/中间缓存资源在什么时候过期, 然后过期后需要重新向服务器获取内容
	// 所以Expires治理利用到了一个过去的时间, 从而可以故意让所有客户端认为过期了，然后不使用缓存
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

	// Last-Modified: 可选，表示资源最后修改时间（用于协商缓存）
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Cors 是一个Gin的中间件, 用来设置options 请求的返回头, 然后退出中间件链, 并结束请求(浏览器跨域设置)
// 实现一个返回头的设置
func Cors(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		// 允许任意来源访问 "*" 则表示不需要携带凭证
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methos", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Header", "authorization, origin, content-type accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}