package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"gongfu/internal/app"
	"gongfu/internal/result"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
	"net/http"
	"strings"
)

func (r UseCase) JSConfig(c *app.Context) result.Result {
	req := struct {
		Uri string `form:"uri" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		return result.Err(err)
	}
	config, err := r.OfficialAccount.GetJs().GetConfig(req.Uri)
	if err != nil {
		return result.Err(fmt.Errorf("get wechat js config: %w", err))
	}
	return result.Ok(gin.H{
		"appID":     config.AppID,
		"timestamp": config.Timestamp,
		"nonceStr":  config.NonceStr,
		"signature": config.Signature,
	})
}

func (r UseCase) Pay(c *app.Context) result.Result {
	type Req struct {
		ActivityId  string `form:"activity_id" binding:"required"`
		Amount      int64  `form:"amount" binding:"required"`
		Description string `form:"description" binding:"required"`
		Attach      string `form:"attach"`
		Username    string `form:"username" binding:"required"`
		Phone       string `form:"phone" binding:"required"`
		Sex         string `form:"sex" binding:"required"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		return result.Err(err)
	}
	if req.Amount <= 0 {
		c.String(http.StatusBadRequest, "amount must > 0")
		return result.Err(nil)
	}
	if c.User.OpenID == nil || *c.User.OpenID == "" {
		c.String(http.StatusBadRequest, "user no openid")
		return result.Err(nil)
	}
	outTradeNo := strings.ReplaceAll(uuid.NewString(), "-", "")
	// create enroll model
	err := r.Store.CreateEnroll(c.Request.Context(), &store.CreateEnrollInput{
		UserId:      c.UserID,
		ActivityId:  req.ActivityId,
		Status:      "pending",
		Amount:      req.Amount,
		OutTradeNo:  outTradeNo,
		Description: req.Description,
		Attach:      req.Attach,
		Username:    req.Username,
		Phone:       req.Phone,
		Sex:         req.Sex,
	})
	if err != nil {
		return result.Err(fmt.Errorf("create enroll failed: %w", err))
	}

	resp, err := r.Wechat.Prepay(c.Request.Context(), &service.PrepayInput{
		OutTradeNo:  outTradeNo,
		Description: req.Description,
		Attach:      req.Attach,
		NotifyUrl:   "https://gongfu.grianchan.com/api/wechat/pay-notify",
		Amount:      req.Amount,
		PayerOpenId: *c.User.OpenID,
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(resp)
}

func (r UseCase) PayNotify(c *app.Context) result.Result {
	transaction := new(payments.Transaction)
	notifyReq, err := r.Wechat.GetNotifyHandler().ParseNotifyRequest(c.Request.Context(), c.Request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		return result.Err(err)
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)
	log.WithContext(c).WithField("transaction", transaction).WithField("notify_req", notifyReq).Info("receive wechat pay notify")

	if transaction.TradeState != nil && *transaction.TradeState == "SUCCESS" {
		// 支付成功处理
		enroll, err := r.Store.GetEnroll(c.Request.Context(), *transaction.OutTradeNo)
		if err != nil {
			return result.Err(fmt.Errorf("get enroll: %w", err))
		}
		if enroll.Status == "success" {
			return result.Ok(nil)
		}
		enroll.Status = "success"
		enroll.TransactionId = *transaction.TransactionId
		if err := r.Store.UpdateEnroll(c.Request.Context(), enroll); err != nil {
			return result.Err(fmt.Errorf("update enroll failed: %w", err))
		}
	}
	return result.Ok(nil)
}

func (r UseCase) GetEnroll(c *app.Context) result.Result {
	activityId := c.Query("activity_id")
	if activityId == "" {
		c.String(http.StatusBadRequest, "activity_id required")
		return result.Err(nil)
	}
	enroll, err := r.Store.GetEnrollByActivityId(c.Request.Context(), activityId)
	if err != nil {
		return result.Err(fmt.Errorf("get enroll: %w", err))
	}
	if enroll == nil {
		return result.Ok(gin.H{
			"id":             0,
			"status":         "",
			"out_trade_no":   "",
			"transaction_id": "",
			"activity_id":    "",
			"amount":         0,
			"created_at":     "",
		})
	}
	return result.Ok(gin.H{
		"id":             enroll.ID,
		"status":         enroll.Status,
		"out_trade_no":   enroll.OutTradeNo,
		"transaction_id": enroll.TransactionId,
		"activity_id":    enroll.ActivityId,
		"amount":         enroll.Amount,
		"created_at":     enroll.CreatedAt.String(),
	})
}
