package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// 定义放置fastgo服务配置的默认目录
	defaultHomeDir = ".fastgo"
	// 指定fasto服务的默认配置文件名
	defaultConfigName = "fg-apiserver.yaml"
)

// 定义需要读取的文件名、环境变量，并将内容读取到viper
func onInitialize() {
	if configFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(configFile)
	} else {
		for _, dir := range searchDirs() {
			// 将dir目录加入到配置文件搜索路径
			viper.AddConfigPath(dir)
		}
		// 设置文件格式为YAML
		viper.SetConfigType("yaml")
		// 配置文件名称
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量并设置前缀
	setupEnvironmentVariables()
	// 读取配置文件
	_ = viper.ReadInConfig()
}

func setupEnvironmentVariables() {
	// 允许viper自动匹配环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀
	viper.SetEnvPrefix("FASTGO")
	// 替换环境变量key中的分隔符
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

// 返回默认的配置文件搜索目录
func searchDirs() []string {
	homeDir, err := os.UserHomeDir()
	// 如果获取用户主目录失败，则打印错误信息并退出
	cobra.CheckErr(err)
	return []string{filepath.Join(homeDir, defaultHomeDir), "."}
}

// filePath 获取默认配置文件的完整路径
func filePath() string {
	home, err := os.UserHomeDir()
	// 如果不能获取用户主目录，则记录错误并返回空路径
	cobra.CheckErr(err)
	return filepath.Join(home, defaultHomeDir, defaultConfigName)
}