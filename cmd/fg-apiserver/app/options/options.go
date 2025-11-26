package options

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type MySQLOptions struct {
	// mapstructure 指的是对应YAML中的同名键
	Addr	string	`json:"addr,omitempty" mapstructure:"addr"`
	Username	string	`json:"username,omitempty" mapstructre:"username"`
	Password	string	`json:"-" mapstructure:"password"`
	Database	string	`json:"database" mapstructure:"database"`
	MaxIdleConnections	int	`json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections,omitempty"`
	MaxOpenConnections	int	`json:"max-open-connections,omitempty" mapstructure:"max-open-connections"`
	MaxConnectionLifeTime	time.Duration	`json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
}

// 用于默认创建一个默认值的
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Addr: "127.0.0.1:3006",
		Username: "onex",
		Password: "onex(#)666",
		Database: "onex",
		MaxIdleConnections: 100,
		MaxOpenConnections: 100,
		MaxConnectionLifeTime: time.Duration(10) *time.Second,
	}
}

type ServerOptions struct {
	// 用来对应YAMl配置文件中的 例如: mysql.addr\mysql.password 等
	MySQLOptions *MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: NewMySQLOptions(),
	}
}

func (o *ServerOptions) Validate() error {
	if o.MySQLOptions.Addr == "" {
		return fmt.Errorf("MySQL server address cannot be empty")
	}
	host, portStr, err := net.SplitHostPort(o.MySQLOptions.Addr)
	if err != nil {
		return fmt.Errorf("Invaild MySQL address format '%s' : %w", o.MySQLOptions.Addr, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("Invalid MySQL port: %s", portStr)
	}
	if host == "" {
		return fmt.Errorf("MySQL hostname cannot be empty")
	}
	if o.MySQLOptions.Username == "" {
		return fmt.Errorf("MySQl username cannot be empty")
	}
	if o.MySQLOptions.Password == "" {
		return fmt.Errorf("MySQL password cannot be empty")
	}
	if o.MySQLOptions.Database == "" {
		return fmt.Errorf("MySQL database cannot be empty")
	}
	if o.MySQLOptions.MaxIdleConnections <= 0 {
		return fmt.Errorf("MySQL max idle connections must be greater than 0")
	}
	if o.MySQLOptions.MaxOpenConnections <= 0 {
		return fmt.Errorf("MySQL max open connections must be greater than 0")
	}
	if o.MySQLOptions.MaxIdleConnections > o.MySQLOptions.MaxOpenConnections {
		return fmt.Errorf("MySQL max idle connections cannot be greater than max open conections")
	}
	if o.MySQLOptions.MaxConnectionLifeTime <= 0 {
		return fmt.Errorf("MySQL max connection lifetime must be greater than 0")
	}
	return nil
}