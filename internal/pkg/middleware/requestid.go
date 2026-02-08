package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/RadishXZ/fastgo/internal/pkg/contextx"
	"github.com/RadishXZ/fastgo/internal/pkg/known"
)

// RequestID 是一个 Gin 中间件, 目的是给每一个 HTTP 请求都发放一个 RequestID
// 作用主要是如果系统报错了, 我们只需要去日志里搜这个 ID 就能看到这个请求从进来到出去发生的所有事情
// RequestID 用来在每一个 HTTP 请求的 context, response 中注入 `x-request-id` 键值对
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 `x-request-id` , 看看请求头里面有没有 x-request-id , 如果不存在的话就获取一个新的 uuid
		requestID := c.Request.Header.Get(known.XRequestID)
		
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 生成一个新的 Context , 并且将 RequestID 保存到 context.Context 中, 以便后续程序使用
		// c.Request.Context 拿到当前原本的上下文
		// contextx.WitchRequestID 跳过了原本的c.Request.Context 这个空请求
		// 整个动作流程: 跳过了原本的空背包, 往里面放一个键值对 Key 是 requestIDKey , Value 是uuid
		// 结果: 返回了一个新的 context(ctx) , 注意: Go 的 Context 是不可变的, 所以想要修改它, 只能生成一个新的
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)

		// 把这个新的 Context 塞入回 HTTP 请求里, 替换掉 HTTP 请求中那个旧的空背包
		// c.Request 对象是指 Go 标准库的 http.Request 结构体
		// WithContext 相当于是拿出了一个新的单(Request_new), 把旧的单(c.Request)上URL、Header、Body 指针等所有信息(浅拷贝)抄过来
		// 替换 Context 的时候, 不会抄旧的 Context_A 这一个, 而是把刚才第一步生成的 Context_B 填进去
		c.Request = c.Request.WithContext(ctx)

		// 将 RequestID 保存到 HTTP 返回头中, Header 的键为 `x-request-id`
		c.Writer.Header().Set(known.XRequestID, requestID)

		// 继续处理请求
		c.Next()
	}
}