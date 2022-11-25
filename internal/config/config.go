package config

import "fmt"

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
	AppID     string
	AppSecret string
}
