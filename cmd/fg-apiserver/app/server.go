package app

import (
	"fmt"

	"github.com/spf13/cobra"
)


func NewFastGOCommand() *cobra.Command {
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
			fmt.Println("Hello FastGo!")
			return nil	
		},
		// 设置命令运行的参数检查，不需要指定命令行参数
		Args: cobra.NoArgs,
	}
	return cmd
}