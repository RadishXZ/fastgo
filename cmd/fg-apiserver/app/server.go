package app

import (
	"encoding/json"
	"fmt"

	"github.com/RadishXZ/fastgo/cmd/fg-apiserver/app/options"
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
			// 将viper中的配置解析到opts变量中
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}
			// 对命令进行校验
			if err := opts.Validate(); err != nil {
				return err
			}
			fmt.Printf("Read MySQl host from Viper: %s\n\n", viper.GetString("mysql.host"))

			jsonData, _ :=json.MarshalIndent(opts, "", " ")
			fmt.Println(string(jsonData))
			return nil	
		},
		// 设置命令运行的参数检查，不需要指定命令行参数
		Args: cobra.NoArgs,
	}
	cobra.OnInitialize(onInitialize)

	// 为命令行程序添加 --config/c 命令行, 用于检查yaml文件的config设置
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the fg-apiserver configuration file.")

	return cmd
}