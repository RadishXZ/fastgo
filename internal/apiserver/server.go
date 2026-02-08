package apiserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	slog.Info("Start to listening the incoming requests on http address", "addr", s.cfg.Addr)
	go func ()  {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	// 设置一个监听信号, 也就是在 channel 中设置一个 Signal通知程序关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 让程序阻塞等待, 等待从 quit channel  中接收到信号
	// quit 收到 SIGTERM 信号后, 程序会解除阻塞状态, 并调用 *http.Server类型实例的 Shutdown方法优雅关停服务器
	<-quit

	slog.Info("Shutting down server ...")

	// 优雅关闭服务
	// 通过context.WithTimeout 设定了10秒, 意思是超过10秒就会强制关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// 先关闭依赖的服务, 在关闭被依赖的服务
	// 10秒内优雅关闭服务 (将未处理完的请求处理完再关闭), 超过 10 秒就超时退出
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Insecure Server forced to shutdown", "err", err)
		return err
	}

	slog.Info("Server exited")

	return nil
}