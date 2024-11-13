package client

import (
	"context"
	"fmt"
	wechat2 "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"gongfu/internal/config"
)

type OfficialAccount struct {
	*officialaccount.OfficialAccount
	subscribeTemplateId string
}

func NewOfficialAccount(conf *config.Config) (*OfficialAccount, error) {
	if conf.WeChat.SubscribeTemplateId == "" {
		return nil, fmt.Errorf("subscribe template id is empty")
	}
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
	return &OfficialAccount{wc.GetOfficialAccount(cfg), conf.WeChat.SubscribeTemplateId}, nil
}

type SubscribeSendInput struct {
	OpenId        string
	CourseName    string
	CourseTime    string
	CourseAddress string
	TeacherName   string
	ToolTip       string
}

func (a OfficialAccount) SubscribeSend(in *SubscribeSendInput) error {
	return a.GetSubscribe().Send(&message.SubscribeMessage{
		ToUser:     in.OpenId,
		TemplateID: a.subscribeTemplateId,
		Page:       "",
		Data: map[string]*message.SubscribeDataItem{
			// 课程名称
			"name1": {
				Value: in.CourseName,
			},
			// 课程时间
			"time2": {
				Value: in.CourseTime,
			},
			// 上课地点
			"thing3": {
				Value: in.CourseAddress,
			},
			// 温馨提示
			"thing4": {
				Value: in.ToolTip,
			},
			// 授课老师
			"thing8": {
				Value: in.TeacherName,
			},
		},
	})
}
