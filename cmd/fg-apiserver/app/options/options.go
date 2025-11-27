package options

import (
	"github.com/RadishXZ/fastgo/internal/apiserver"
	genericoptions "github.com/RadishXZ/fastgo/pkg/options"
)

// 调用 pkg/options/mysql_options.go 嵌入到Server配置中
type ServerOptions struct {
	// 用来对应YAMl配置文件中的 例如: mysql.addr\mysql.password 等
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql“ mapstructure:"mysql"`	
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: genericoptions.NewMySQLOptions(),
	}
}

func (o *ServerOptions) Validate() error {
	if err := o.MySQLOptions.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		MySQLOptions: o.MySQLOptions,
	}, nil
}