package service

import (
	"context"
	wechat2 "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"gongfu/internal/config"
)

func NewOfficialAccount(conf *config.Config) *officialaccount.OfficialAccount {
	wc := wechat2.NewWechat()
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
