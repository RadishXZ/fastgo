## fastgo 项目
云原生 Go 练习实战

### 总体流程顺序

    cmd/fg-apiserver/main.go --> 
    cmd/fg-apiserver/app.NewFastGOCommand() --> cobra执行 --> onInitialize() (viper读取配置) --> run(opts) --> viper.Unmarshal(opts) --> opts.Validate() --> opts.Config() --> cfg.NewServer() --> server.Run()



### 应用入口

    cmd/fg-apiserver/main.go
    程序运行的主入口, 调用cobra命令触发启动流程 --->



### 应用框架代码

    cmd/fg-apiserver/app/options.go
    聚合CLI配置(ServerOptions) ，通过(MySQLOptions指针) 提供Validate校验、Config配置项，从而生成 apiserver.Config --->

    cmd/fg-apiserver/app/server.go
    构建cobra命令, 并在RunE调用run(opts), run负责把viper配置反序列化struct: viper.Unmarshal(opts)并且触发后续启动流程 --->



### 读取配置文件、viper初始化

    cmd/fg-apiserver/app/config.go
    把外部的 yaml文件、环境变量、命令行 统一读取到程序运行的配置中心(viper), 并为后续的viper.Unmarshal(opts)提供已加载的键值



### 运行时代码

    internal/apiserver/server.go
    实际上算是个运行时的汇总入口
        --- 接收 options 已验证的 apiserver.Config, 实际关系在(pkg/options/mysql_options.go) 中 MySQLOptions 的结构体去
        --- 然后通过Server.Run()读取以上配置信息创建资源(比如调用MySQLOptions.NewDB())、启动后台服务并管理生命周期 --->



### 作为共享配置信息代码

    pkg/options/mysql_options.go
    定义MySQLOptions (Validate、DSN、NewDB) 字段，提供给ServerOptions所需要的配置信息、校验、执行动作 --->



### 通过viper包解析yaml格式文件读取配置信息

    configs/fg-apiserver.yaml
    viper读取里面的yaml内容，然后映射到(cmd/fg-apiserver/app/options/options.go)的ServerOptions (依赖mapstructure指定的内容)


