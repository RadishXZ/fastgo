package version

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

// 定义了 --version 参数的三种状态
type versionValue int

const (
	VersionNotSet versionValue = 0	// 0 = 不打印, 继续运行程序
	VersionEnabled versionValue = 1	// 1 = 打印简短版本 (v0.0.1)
	VersionRaw versionValue = 2	// 2 = 打印详细版本 (完整表格)
)

const strRawVersion string = "raw"

// 接口实现, 让pflag能识别你的自定义类型是布尔
func (v *versionValue) IsBoolFlag() bool {
	return true
}

// 接口实现, 让pflag能识别的自定义类型是any
func (v *versionValue) Get() any {
	return *v
}

// 用户输入了什么, 就在这里处理
func(v *versionValue) Set(s string) error {
	if s == strRawVersion { // 如果是 --version=raw
		*v = VersionRaw	// 输出为详细版本信息
		return nil
	}
	// 否则当做布尔值处理
	boolVal, err := strconv.ParseBool(s)
	if boolVal {	// 如果是true
		*v = VersionEnabled	// 设置为 简短版本
	} else {	// 如果是false
		*v = VersionNotSet	// 设置为不打印
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion	// 显示为 raw
	}
	return fmt.Sprintf("%v", bool(*v == VersionEnabled))	// 显示true或者false
}

func (v *versionValue) Type() string {
	// 在帮助信息中显示这个 flag 类型
	return "version"
}
// 作用是把 versionValue类型注册为命令行参数flag, 让程序支持 --version 参数, 并且能处理true\false\=raw 三种情况
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value // 设置初始值
	flag.Var(p, name, usage) // 把这个自定义类型注册到 flag 系统
	flag.Lookup(name).NoOptDefVal = "true" // 当用户写 --version 不带值, 自动用true作为值
}

// 注册一个新的versionValue类型的命令行参数flag, 并返回这个参数的指针
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)	// 创建一个新的versionValue
	VersionVar(p, name, value, usage)	// 注册自定义类型
	return p	// 返回指针, 提供后面使用
}

const versionFlagName = "version"

// 最终 pflag 会解析*versionFlag 并修改
var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit")

// 作用是 程序启动时就创建并注册一个全局的 --version的flag
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// 把这个flag加到你的cobra中
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		// 用户输入指令是详细版本的话
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)	// 立即退出程序
	} else if *versionFlag == VersionEnabled {
		// 用户输入指令是简短版本的话
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)	// 立即退出程序
	}
}