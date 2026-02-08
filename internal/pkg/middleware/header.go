package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 是一个 Gin 的中间件, 用来禁止客户端缓存 HTTP 请求的返回结果.
func NoCache(c *gin.Context) {

	// 告诉浏览器 "这个接口的数据经常变, 不存缓存, 需要每次都得找服务器要最新的请求"
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	
	// 设置HTTP 的 Expires 响应头告诉客户端/中间缓存资源在什么时候过期, 然后过期后需要重新向服务器获取内容
	// 所以Expires治理利用到了一个过去的时间, 从而可以故意让所有客户端认为过期了，然后不使用缓存
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

	// Last-Modified: 可选，表示资源最后修改时间（用于协商缓存）
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	// 把控制权交给下一个中间件或者最终处理函数
	c.Next()
}

// Cors 是一个Gin的中间件, 用来作为安全机制, 因为默认情况下, 浏览器不允许 A 网站前端代码去请求 B 网站的后端接口
// Cors 就是后端告诉浏览器: 我是允许 A 网站来访问我的
// 实现一个返回头的设置
func Cors(c *gin.Context) {

	// 如果是正常的请求, 就会往下走 c.Next()
	if c.Request.Method != "OPTIONS" {
		c.Next()
	
	// 如果是 OPTIONS 请求, 则直接在 Header 里面告诉允许哪些方法, 允许哪些头部信息(GET\POST\PUT ...), 最后通过 c.AborWithStatus 结束
	// 这其实是浏览器的探路请求, 在正式发送数据前, 浏览器会先问一句: 我能连你么
	} else {
		// 允许任意来源访问 "*" 则表示不需要携带凭证
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(200)
	}
}