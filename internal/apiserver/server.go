package apiserver

import (
	"errors"
	"log/slog"
	"net/http"

	mw "github.com/RadishXZ/fastgo/internal/pkg/middleware"
	genericoptions "github.com/RadishXZ/fastgo/pkg/options"
	"github.com/gin-gonic/gin"
)

// Config 配置结构体, 用于存储应用相关的配置
type Config struct {
	MySQLOptions *genericoptions.MySQLOptions
	Addr string
}

// Server 定义一个服务器结构体类型
type Server struct {
	cfg *Config
	srv *http.Server
}

// NewServer 根据配置构建并返回一个包含 Gin 引擎的 http.Server (不启动监听)
func (cfg *Config) NewServer() (*Server, error) {

	// 创建一个Gin引擎
	engine := gin.New()

	// gin.Recovery() 中间件, 用来捕获 panic 并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.RequestID()}
	engine.Use(mws...)
	
	// 注册404, 代码通过 JSON 返回 404
	engine.NoRoute(func (c *gin.Context)  {
		// 调用
		c.JSON(http.StatusNotFound, gin.H{"code": "PageNotFound", "message": "Page not found."})
	})

	// 注册 /healthz Handler处理器, 检查注册健康检查接口, 检查我们任何觉得会影响服务器健康状态的项目
	engine.GET("/healthz", func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 构造标准库的http.Server 把Gin引擎作为请求处理器,并把监听地址设为cfg.Addr
	// 将 gin 引擎包装为标准库 http.Server，便于设置超时、TLS、优雅停机等
	httpsrv := &http.Server{Addr: cfg.Addr, Handler: engine}
	return  &Server{cfg: cfg, srv: httpsrv}, nil
}

// Run 运行应用入口
// 输出信息示例: time=2026-01-28T15:13:10.257+08:00 level=INFO msg="Read MySQL host from config" mysql.addr=127.0.0.1:3306
func (s *Server) Run() error {
	// 运行 HTTP 服务器
	// 打印一条日志, 用来提示 HTTP 服务已经起来, 调用slog记录一条info级别的日志信息, 内容格式以信息+地址呈现
	slog.Info("Read MySQL host from config", "mysql.addr", s.cfg.MySQLOptions.Addr)
	// fmt.Printf("Read MySQL host from config: %s\n", s.cfg.MySQLOptions.Addr)  # 日志功能加入, 已被替代
	
	// 调用s.srv.ListenAndServe方法启动服务器, 当该方法返回错误时, 报错退出
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}