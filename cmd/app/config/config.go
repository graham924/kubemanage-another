package config

// SysConfig 系统配置，全局对象
var SysConfig *Config

// Config 配置对象
type Config struct {
	Default DefaultOptions `mapstructure:"default"`
	Mysql   MysqlOptions   `mapstructure:"mysql"`
	Log     LogConfig      `mapstructure:"log"`
}

// DefaultOptions 默认配置选项
type DefaultOptions struct {
	PodLogTailLine       string `mapstructure:"podLogTailLine"`
	ListenAddr           string `mapstructure:"listenAddr"`
	WebSocketListenAddr  string `mapstructure:"webSocketListenAddr"`
	JWTSecret            string `mapstructure:"JWTSecret"`
	ExpireTime           int64  `mapstructure:"expireTime"`
	KubernetesConfigFile string `mapstructure:"kubernetesConfigFile"`
}

// MysqlOptions mysql配置选项
type MysqlOptions struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Port         string `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxLifetime  int    `mapstructure:"maxLifetime"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
}

// LogConfig 日志配置选项
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
