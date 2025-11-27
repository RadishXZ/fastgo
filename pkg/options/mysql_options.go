package options

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 给 MySQL提供默认值
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

// 用于默认创建一个默认值, 定义是一个返回结构体指针的函数, 主要都是些默认值为主
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

// 校验填写逻辑是否规范
func (o *MySQLOptions) Validate() error {
	if o.Addr == "" {
		return fmt.Errorf("MySQL server address cannot be empty")
	}
	host, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("Invaild MySQL address format '%s' : %w", o.Addr, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("Invalid MySQL port: %s", portStr)
	}
	if host == "" {
		return fmt.Errorf("MySQL hostname cannot be empty")
	}
	if o.Username == "" {
		return fmt.Errorf("MySQl username cannot be empty")
	}
	if o.Password == "" {
		return fmt.Errorf("MySQL password cannot be empty")
	}
	if o.Database == "" {
		return fmt.Errorf("MySQL database cannot be empty")
	}
	if o.MaxIdleConnections <= 0 {
		return fmt.Errorf("MySQL max idle connections must be greater than 0")
	}
	if o.MaxOpenConnections <= 0 {
		return fmt.Errorf("MySQL max open connections must be greater than 0")
	}
	if o.MaxIdleConnections > o.MaxOpenConnections {
		return fmt.Errorf("MySQL max idle connections cannot be greater than max open conections")
	}
	if o.MaxConnectionLifeTime <= 0 {
		return fmt.Errorf("MySQL max connection lifetime must be greater than 0")
	}
	return nil
}

func (o *MySQLOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&lov=%s`,
	o.Username,
	o.Password,
	o.Addr,
	o.Database,
	true,
	"Local")
}

func (o *MySQLOptions) NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(o.DSN()), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(o.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(o.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(o.MaxIdleConnections)

	return db, nil
}