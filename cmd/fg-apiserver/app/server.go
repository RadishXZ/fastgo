package app

import (
	"io"
	"log/slog"
	"os"

	"github.com/RadishXZ/fastgo/cmd/fg-apiserver/app/options"
	"github.com/RadishXZ/fastgo/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string


func NewFastGOCommand() *cobra.Command {
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		// 指定命令的名字，改名字会出现在帮助信息中
		Use: "fg-apiserver",
		// 命令的简短描述
		Short: "A very lightweight full go project",
		Long: `A very lightweight full go project, designed to help beginners quickly
		learn Go project development.`,
		SilenceUsage: true,	// 设置true的话，出错时不打印帮助信息，直接输出错误信息
		// 指定调用cmd.Excute()时执行Run函数
		RunE: func (cmd *cobra.Command, args []string) error  {
			return run(opts)
		},
		// 设置命令运行的参数检查，不需要指定命令行参数
		Args: cobra.NoArgs,
	}
	cobra.OnInitialize(onInitialize)

	// 为命令行程序添加 --config/c 命令, 用于检查yaml文件的config设置
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the fg-apiserver configuration file.")

	// 为命令行程序添加 --version 命令
	version.AddFlags(cmd.PersistentFlags())
	return cmd
}

// run 是主运行逻辑, 负责初始化日志、解析配置、校验选项并启动服务器
func run (opts *options.ServerOptions) error {
	// viper.Unmarshal 会把配置文件的值自动填充到opts结构体中
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}
	if err := opts.Validate(); err != nil {
		return err
	}
	
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	server, err := cfg.NewServer()
	if err != nil {
		return err
	}

	// 如果传入 --version 则打印版本信息并退出
	version.PrintAndExitIfRequested()

	initLog()

	return server.Run()
}

// 初始化全局日志实例, 默认全局logger
func initLog() {
	// 获取日志配置
	format := viper.GetString("log.format")	// 日志格式支持 json text
	level := viper.GetString("log.level")	// 日志级别对应映射 debug->LevelDebug info->LevelInfo warn->LevelWarn
	output := viper.GetString("log.output")	// 日志输出路径支持 标准输出stdout和文件

	// 转换日志级别, debug info warn
	var slevel slog.Level
	switch level {
	case "debug":
		slevel = slog.LevelDebug
	case "info":
		slevel = slog.LevelInfo
	case "warn":
		slevel = slog.LevelWarn
	default:
		slevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: slevel}

	var w io.Writer
	var err error

	// 转换日志输出路径
	switch output {
	case "":
		w = os.Stdout
	case "stdout":
		w = os.Stdout
	default:
		// os.OpenFile: 打开或者创建一个文件并返回一个记录, 并实现io.Writer写入日志
		// os.O_CREATE: 文件不存在的话则创建、os.O_WRONLY: 只写、os.O_APPEND: 写入的时候追加到文件末尾
		// 权限0666 代表: rw-rw-rw-
		w, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 转换日志格式
	if err != nil {
		return
	}

	// 声明变量handler 类型是slog.Handler(日志处理器接口), Handler负责把日志格式化并写入底层输出(stdout\文件)
	var handler slog.Handler
	switch format {
	// 日志以json编码格式保存
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	// 日志格式以可读写文本保存
	case "text":
		handler = slog.NewTextHandler(w, opts)
	// 默认以json编码保存
	default:
		handler = slog.NewJSONHandler(w, opts)
	}

	// 负责把时间格式化并写入 w , 并构造一个新的*slog.logger实例, logger通过info、error、debug等方法记录日志
	slog.SetDefault(slog.New(handler))
}