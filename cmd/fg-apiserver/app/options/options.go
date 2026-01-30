package options

import (
	"fmt"
	"net"
	"strconv"

	"github.com/RadishXZ/fastgo/internal/apiserver"
	genericoptions "github.com/RadishXZ/fastgo/pkg/options"
)

// 调用 pkg/options/mysql_options.go 嵌入到Server配置中
type ServerOptions struct {
	// 用来对应YAMl配置文件中的 例如: mysql.addr\mysql.password 等
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql“ mapstructure:"mysql"`
	Addr	string	`json:"addr" mapstructure:"addr"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: genericoptions.NewMySQLOptions(),
		Addr: "0.0.0.0:6666",
	}
}

// Validate 校验 ServerOptions 仲的选项是否合法
func (o *ServerOptions) Validate() error {
	if err := o.MySQLOptions.Validate(); err != nil {
		return err
	}

	if o.Addr == "" {
		return fmt.Errorf("server address cannot be empty")
	}

	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invaild server address format '%s': %w", o.Addr, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invaild server port: %s", portStr)
	}
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		MySQLOptions: o.MySQLOptions,
		Addr: o.Addr,
	}, nil
}