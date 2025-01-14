package global

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Mode          string // debug | release
		Port          string
		DbType        string // mysql | sqlite
		DbAutoMigrate bool   // whether to automatically migrate database table structures
		DbLogMode     string
	}

	Log struct {
		Level     string // debug | info | warn | error
		Prefix    string
		Format    string // text | json
		Directory string
	}
	JWT struct {
		Secret string
		Expire int64
		Issuer string
	}
	Mysql struct {
		Host     string // server address
		Port     string
		Config   string
		Dbname   string // databae name
		Username string // database user name
		Password string
	}
	SQLite struct {
		Dsn string
	}
	Redis struct {
		DB       int
		Addr     string
		Password string
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
	Email struct {
		From     string // From, the email address you want to send it to
		Host     string // Serveraddress,such as smtp qq.com Go to the mailbox you want to send an email to view its smtp protocol
		Port     int    // 前往要发邮件的邮箱查看其 smtp 协议端口, 大多为 465
		SmtpPass string // email secret key
		SmtpUser string
	}
	Captcha struct {
		SendEmail  bool // 是否通过邮箱发送验证码
		ExpireTime int  // 过期时间
	}
	Upload struct {
		// Size
		OssType   string // local | qiniu
		Path      string // 本地访问路径
		StorePath string // 本地文件存储路径
	}
	Qiniu struct {
		ImPath        string // 外部链接
		Zone          string // 存储区域
		Bucket        string // 空间名称
		AccessKey     string // secret key AK
		SecretKey     string // secret key SK
		UseHTTPS      bool   // whether use https
		UseCdnDomains bool   // 上传是否使用CDS 上传加速
	}
}

var Conf *Config

func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

// 从指定路径读取配置文件
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()                                   // 允许使用环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // SERVER_APPMODE => SERVER>APPMODE

	if err := v.ReadConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}

	log.Println("配置文件加载成功:", path)
	return Conf
}

// 数据库类型
func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

// 数据库连接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		conf := Conf.Mysql
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
		)
	case "sqlite":
		return Conf.SQLite.Dsn
	// 默认使用 sqlite, 并且使用内存数据库
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.SQLite.Dsn == "" {
			Conf.SQLite.Dsn = "file::memory:"
		}
		return Conf.SQLite.Dsn
	}
}
