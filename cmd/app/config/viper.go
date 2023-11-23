package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var configObj = &Config{
	Default: DefaultOptions{
		ListenAddr:          ":6180",
		WebSocketListenAddr: "",
		JWTSecret:           "kubemanage",
		ExpireTime:          10,
	},
	Mysql: MysqlOptions{
		Host:         "192.168.245.100",
		Port:         "3306",
		User:         "root",
		Password:     "zgy123.com",
		Name:         "kubemanage",
		MaxOpenConns: 100,
		MaxLifetime:  20,
		MaxIdleConns: 10,
	},
	Log: LogConfig{
		Level:      "debug",
		Filename:   "kubemanage.log",
		MaxSize:    200,
		MaxAge:     30,
		MaxBackups: 7,
	},
}

// Binding 解析外部的配置文件，默认是 ./config.yaml
func Binding(filePath string) error {
	v := viper.New()
	// TODO 镜像里，config的路径总有问题，先使用配置对象规避掉
	//v.SetConfigFile(filePath)
	//if err := v.ReadInConfig(); err != nil {
	//	return err
	//}
	//// 把读取到的配置信息反序列化到 SysConfig 变量中
	//if err := v.Unmarshal(&SysConfig); err != nil {
	//	return fmt.Errorf("config Unmarshal failed, err:%v\n", err)
	//}
	SysConfig = configObj
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed,sys config reload")
		if err := viper.Unmarshal(&SysConfig); err != nil {
			fmt.Printf("config file changed,viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return nil
}
