package app

import (
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

	return server.Run()
}