package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type SMS struct {
	AppID      string
	TemplateID string
	SignName   string
}

type COS struct {
	BucketURL string
}

type Tencent struct {
	SecretId  string
	SecretKey string
	SMS       SMS
	COS       COS
}

type Redis struct {
	Addr string
	DB   int
}

type Token struct {
	Secret string
}

type DB struct {
	Addr     string
	User     string
	Password string
	DB       string
}

func (d DB) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Addr, d.DB)
}

type Logger struct {
	Level  string
	Format string
}

type Config struct {
	WeChat  WeChat
	Tencent Tencent
	Redis   Redis
	Token   Token
	DB      DB
	// web 页面的本地相对路径
	WebPath string
	// 当前环境 production、dev、testing
	Env    string
	Logger Logger
}

func (c Config) IsProd() bool {
	return c.Env == "production"
}

type WeChat struct {
	AppID                      string
	AppSecret                  string
	MchID                      string // 商户号
	MchCertificateSerialNumber string // 商户证书序列号
	MchAPIv3Key                string // 商户APIv3密钥
	PrivateCertPath            string
	SubscribeTemplateId        string // 订阅消息模板id
}

func NewConfig(file string) (*Config, error) {
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &conf, nil
}
