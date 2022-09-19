package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"strconv"
	"time"
)

type AuthCode interface {
	Send(phone string) (code string, err error)
	Valid(phone string, code string) (bool, error)
}

type authCode struct {
	Client      *sms.Client
	Redis       redis.Cmdable
	SmsSdkAppId string
	SignName    string
	TemplateId  string
	Period      time.Duration
	CodeLength  int
	SendLimit   int
	Rand        *rand.Rand
}

func NewAuthCode(
	client *sms.Client,
	redis redis.Cmdable,
	smsSdkAppId string,
	signName string,
	templateId string,
	period time.Duration,
	codeLength int,
) *authCode {
	return &authCode{
		Client:      client,
		Redis:       redis,
		SmsSdkAppId: smsSdkAppId,
		SignName:    signName,
		TemplateId:  templateId,
		Period:      period,
		CodeLength:  codeLength,
		SendLimit:   10,
		Rand:        rand.New(rand.NewSource(time.Now().Unix())),
	}
}

type SendErr struct {
	Err      error
	Response *sms.SendSmsResponse
}

func (e SendErr) IsError() bool {
	if e.Err != nil {
		return true
	}
	if e.Response != nil && e.Response.Response != nil {
		if len(e.Response.Response.SendStatusSet) > 0 {
			code := e.Response.Response.SendStatusSet[0].Code
			if code != nil && *code != "Ok" {
				return true
			}
		}
	}
	return false
}

func (e SendErr) GetSendStatus() *sms.SendStatus {
	if e.Response == nil {
		return nil
	}
	if e.Response.Response == nil {
		return nil
	}
	if len(e.Response.Response.SendStatusSet) == 0 {
		return nil
	}
	return e.Response.Response.SendStatusSet[0]
}

func (e SendErr) ToError() error {
	if !e.IsError() {
		return nil
	}
	status := e.GetSendStatus()
	if status == nil {
		return fmt.Errorf("%w", e.Err)
	} else {
		return fmt.Errorf("err=%s, requestID=%s, phone=%s, serialNo=%s, code=%s, message=%s",
			e.Err,
			*e.Response.Response.RequestId,
			*status.PhoneNumber,
			*status.SerialNo,
			*status.Code,
			*status.Message,
		)
	}
}

func (a authCode) Send(phone string) (code string, err error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if isLimiting, err := a.isLimiting(c, phone); err != nil {
		return "", fmt.Errorf("check limiting: %w", err)
	} else if isLimiting {
		return "", fmt.Errorf("in limiting")
	}
	code = a.newCode(a.CodeLength)
	period := strconv.Itoa(int(a.Period.Minutes()))
	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = []*string{&phone}
	req.SmsSdkAppId = &a.SmsSdkAppId
	req.TemplateId = &a.TemplateId
	req.SignName = &a.SignName
	req.TemplateParamSet = []*string{&code, &period}
	resp, err := a.Client.SendSmsWithContext(c, req)
	sendErr := SendErr{
		Err:      err,
		Response: resp,
	}
	if err := sendErr.ToError(); err != nil {
		return "", fmt.Errorf("send sms: %w", err)
	}

	if err := a.Redis.Set(c, fmt.Sprintf("sms_valid_code_%s", phone), code, a.Period).Err(); err != nil {
		return "", fmt.Errorf("cache code: %w", err)
	}
	_ = a.addLimiting(c, phone)
	return code, nil
}

func (a authCode) Valid(phone string, code string) (bool, error) {
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if code2, err := a.Redis.Get(c, fmt.Sprintf("sms_valid_code_%s", phone)).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("get cache code: %w", err)
	} else {
		if code2 != code {
			return false, nil
		} else {
			return true, nil
		}
	}
}

func (a authCode) newCode(length int) string {
	if length < 1 || length > 10 {
		length = 4
	}
	code := ""
	for i := 0; i < length; i++ {
		code += strconv.Itoa(a.Rand.Intn(10))
	}
	return code
}

func (a authCode) addLimiting(ctx context.Context, phone string) error {
	return a.Redis.Incr(ctx, fmt.Sprintf("limiting_%s_%s", phone, time.Now().Format("2006050415"))).Err()
}

func (a authCode) isLimiting(ctx context.Context, phone string) (bool, error) {
	v, err := a.Redis.Get(ctx, fmt.Sprintf("limiting_%s_%s", phone, time.Now().Format("2006050415"))).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if i, err := strconv.Atoi(v); err != nil {
		return false, err
	} else {
		return i >= 10, nil
	}
}
