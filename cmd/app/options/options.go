package options

import (
	"fmt"
	localLog "log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/noovertime7/kubemanage/cmd/app/config"
	"github.com/noovertime7/kubemanage/dao"
	"github.com/noovertime7/kubemanage/pkg"
	log "github.com/noovertime7/kubemanage/pkg/logger"
	"github.com/noovertime7/kubemanage/pkg/source"
)

const (
	// defaultConfigFile 默认配置文件路径
	defaultConfigFile = "./config.yaml"
)

// Options 命令行选项对象
type Options struct {
	// GinEngine gin引擎对象
	GinEngine *gin.Engine
	// The default values.
	// DB 数据库对象
	DB *gorm.DB
	// ShareDaoFactory 数据库抽象工厂 包含所有数据操作接口
	Factory dao.ShareDaoFactory // 数据库接口
	// 配置文件路径
	ConfigFile string
}

// NewOptions 新建一个命令行选项对象
func NewOptions() (*Options, error) {
	return &Options{
		// 设置默认的配置文件路径
		ConfigFile: defaultConfigFile,
	}, nil
}

// BindFlags 绑定配置文件路径的命令行参数
func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "configFile", "", "The location of the kubemanage configuration file")
}

// Complete completes all the required options
func (o *Options) Complete() error {
	// 配置文件优先级，从都到高为：默认配置，环境变量，命令行
	if len(o.ConfigFile) == 0 {
		// Try to read config file path from env.
		if cfgFile := os.Getenv("KUBEMANAGE-CONFIG"); cfgFile != "" {
			o.ConfigFile = cfgFile
		} else {
			o.ConfigFile = defaultConfigFile
		}
	}

	// 解析配置文件
	if err := config.Binding(o.ConfigFile); err != nil {
		return err
	}

	// 初始化默认 api 路由
	o.GinEngine = gin.Default()

	// 注册依赖组件
	if err := o.register(); err != nil {
		return err
	}
	return nil
}

// InitDB 初始化DB
func (o *Options) InitDB() error {
	initDbService := source.NewInitDBService(o.DB)
	return initDbService.InitDB()
}

// register 注册依赖组件
func (o *Options) register() error {
	// 注册日志组件
	if err := o.registerLogger(); err != nil {
		return err
	}
	// 注册数据库组件
	if err := o.registerDatabase(); err != nil {
		return err
	}
	return nil
}

// registerLogger 注册日志组件
func (o *Options) registerLogger() error {
	return log.InitLogger()
}

// registerDatabase 注册数据库组件
func (o *Options) registerDatabase() error {
	newLogger := logger.New(
		localLog.New(os.Stdout, "\r\n", localLog.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	// 取出sql配置
	sqlConfig := config.SysConfig.Mysql
	// 拼接mysql的dsn链接（TODO）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		sqlConfig.User,
		sqlConfig.Password,
		sqlConfig.Host,
		sqlConfig.Port,
		sqlConfig.Name)
	var err error
	if o.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	}); err != nil {
		return err
	}
	// 设置数据库连接池，使用的是options中创建好的gorm对象
	sqlDB, err := o.DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(config.SysConfig.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.SysConfig.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.SysConfig.Mysql.MaxLifetime) * time.Second)
	o.Factory = dao.NewShareDaoFactory(o.DB)
	return nil
}

func (o *Options) registerJwt() {
	pkg.RegisterJwt(config.SysConfig.Default.JWTSecret)
}
