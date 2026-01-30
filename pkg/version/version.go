package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	gitVersion	= "v0.0.0-master+$Format:%H$"
	gitCommit = "$Format:%H$"
	gitTreeState = ""
	buildDate = "1970-01-01T00:00:00Z"
)

type Info struct {
	GitVersion string `json:"gitVersion"`	// 版本号
	GitCommit string `json:"gitCommit"`	//	git 提交哈希
	GitTreeState string `json:"gitTreeState"`	//	源码状态, clean没有改动, dirty表示有
	BuildDate string `json:"buildDate"`	//	编译时间
	GoVersion string `json:"goVersion"` 	//	Go语言版本
	Compiler string `json:"compiler"`	//	编译器, 一般是gc
	Platform string `json:"platform"`	//	平台信息, darwin/amd64, linux/arm64, risv

}

//	返回版本号, 最简单展示构建信息
func (info Info) String() string {
	return info.GitVersion
}

//	把整个结构体转成JSON字符串
func (info Info) ToJSON() string {
	// json.Marshal 可以把Go常见的结构体转换成JSON格式
	s, _ := json.Marshal(info)
	return string(s)
}

// 用表格形式打印出来, 展示格式化的版本信息
func (info Info) Text() string {
	// 创建一个新的表格 (空表)
	table := uitable.New()
	// 把指定列设置为向右对齐, 0是第一列, 也可以传入多个参数 (0, 2) 表示第1和3列右对齐
	table.RightAlign(0)
	// 限制单列最大字符宽度, 多出来的会自动把内容拆成多行自适应
	table.MaxColWidth = 80
	// 设置烈玉列之间用作分隔的字符串, " "表示空格， 也可以指定其他的作为分隔符, 例如 " | " 
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.String()
}

// 收集所有信息, 返回详细的代码库版本信息
func Get() Info {
	return Info {
		GitVersion: gitVersion,
		GitCommit: gitCommit,
		GitTreeState: gitTreeState,
		BuildDate: buildDate,

		// 调用 runtime来显示程序中运行的一些环境信息 
		GoVersion: runtime.Version(),
		Compiler: runtime.Compiler,
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

