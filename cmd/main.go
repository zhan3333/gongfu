package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gongfu/internal/config"
	http2 "gongfu/internal/http"
	"gongfu/internal/http/controller"
	"gongfu/internal/http/middlewares"
	"gongfu/internal/model"
	"gongfu/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"time"
)

func main() {
	viper.SetConfigFile("config/production.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic("read config failed: " + err.Error())
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic("unmarshal config failed: " + err.Error())
	}
	initEnv(conf.Env)
	if err := initLogger(&conf.Logger); err != nil {
		panic(err)
	}
	r := gin.Default()
	oa := getOfficialAccount(&conf)
	redis2, err := getRedis(&conf.Redis)
	if err != nil {
		panic("init redis: " + err.Error())
	}
	authCode := getAuthCode(&conf.Tencent, redis2)
	token := service.NewToken(conf.Token.Secret)
	storageService := getStorageService(&conf.Tencent)

	store, err := getStore(&conf.DB)
	if err != nil {
		panic(err)
	}
	control := controller.NewController(&conf, oa, authCode, token, store, storageService)
	middleware := middlewares.NewMiddlewares(token, store)

	http2.NewRoute(&conf, oa, control, middleware).Route(r)
	if err := r.Run("0.0.0.0:9003"); err != nil {
		fmt.Println(err)
	}
}

func getStorageService(conf *config.Tencent) service.Storage {
	u, _ := url.Parse(conf.COS.BucketURL)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})
	return service.NewStorage(client, conf.SecretId, conf.SecretKey)
}

func initEnv(env string) {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if env == "test" {
		gin.SetMode(gin.TestMode)
	}
}

func initLogger(conf *config.Logger) error {
	lv, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lv)
	if conf.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}

func getStore(conf *config.DB) (service.Store, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.CheckIn{},
		&model.CheckInDay{},
		&model.CheckInCount{},
		&model.Coach{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return service.NewDBStore(db), nil
}

func getAuthCode(conf *config.Tencent, redis redis.Cmdable) service.AuthCode {
	credential := common.NewCredential(conf.SecretId, conf.SecretKey)
	client, _ := sms.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	return service.NewAuthCode(client, redis, conf.SMS.AppID, conf.SMS.SignName, conf.SMS.TemplateID, 5*time.Minute, 4)
}

func getRedis(conf *config.Redis) (redis.Cmdable, error) {
	client := redis.NewClient(&redis.Options{
		Addr: conf.Addr,
		DB:   conf.DB,
	})
	return client, client.Ping(context.Background()).Err()
}

func getOfficialAccount(conf *config.Config) *officialaccount.OfficialAccount {
	wc := wechat.NewWechat()
	//设置全局cache，也可以单独为每个操作实例设置
	redisOpts := &cache.RedisOpts{
		Host:     conf.Redis.Addr,
		Database: conf.Redis.DB,
	}
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	wc.SetCache(redisCache)
	cfg := &offConfig.Config{
		AppID:     conf.WeChat.AppID,
		AppSecret: conf.WeChat.AppSecret,
		//Token:     "xxx",
		//EncodingAESKey: "xxxx",
		//Cache: redisCache, //也可以单独设置
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	return officialAccount
}
