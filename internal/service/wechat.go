package service

import (
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"gongfu/internal/config"
)

type wechat struct {
	client        *core.Client
	AppId         string
	MchId         string
	JsapiService  jsapi.JsapiApiService
	NotifyHandler *notify.Handler
}

func (w *wechat) SubscribeSend(ctx context.Context, in *SubscriptionSendInput) error {
	return nil
}

type SubscriptionSendInput struct {
}

type Wechat interface {
	GetNotifyHandler() *notify.Handler
	Prepay(ctx context.Context, in *PrepayInput) (*jsapi.PrepayWithRequestPaymentResponse, error)
	SubscribeSend(ctx context.Context, in *SubscriptionSendInput) error
}

func NewWechat(config *config.Config) (Wechat, error) {
	var (
		mchID                      string = config.WeChat.MchID                      // 商户号
		mchCertificateSerialNumber string = config.WeChat.MchCertificateSerialNumber // 商户证书序列号
		mchAPIv3Key                string = config.WeChat.MchAPIv3Key                // 商户APIv3密钥
	)

	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(config.WeChat.PrivateCertPath)
	if err != nil {
		return nil, fmt.Errorf("load merchant private key error: %w", err)
	}

	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("new wechat pay client err:%w", err)
	}

	// init notify handler
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, mchPrivateKey, mchCertificateSerialNumber, mchID, mchAPIv3Key)
	if err != nil {
		return nil, fmt.Errorf("register downloader: %w", err)
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(mchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		return nil, fmt.Errorf("new notify handler: %w", err)
	}

	return &wechat{
		client:        client,
		AppId:         config.WeChat.AppID,
		MchId:         mchID,
		JsapiService:  jsapi.JsapiApiService{Client: client},
		NotifyHandler: handler,
	}, nil
}

type PrepayInput struct {
	OutTradeNo  string
	Description string
	Attach      string
	NotifyUrl   string
	Amount      int64
	PayerOpenId string
}

func (w *wechat) Prepay(ctx context.Context, in *PrepayInput) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	resp, _, err := w.JsapiService.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(w.AppId),
			Mchid:       core.String(w.MchId),
			Description: core.String(in.Description),
			OutTradeNo:  core.String(in.OutTradeNo),
			Attach:      core.String(in.Attach),
			NotifyUrl:   core.String(in.NotifyUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(in.Amount),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(in.PayerOpenId),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("call prepay api error: %w", err)
	}
	return resp, nil
}

func (w *wechat) GetNotifyHandler() *notify.Handler {
	return w.NotifyHandler
}
