package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/pkg/util"
	"net/http"
	"strings"
)

var wechatLoginURL = "/web/login"

// WeChatLogin 公众号登录
// /login 时会 302 到微信授权登录页面, 用户确认授权后返回到 /login?code={} 地址
// /login?code={} 时处理 code 为 accessToken 获取微信用户信息
// 获取微信用户信息后，根据 openID 查找用户并生成 JWT token
// 生成 JWT token 后 302 到 /me?accessToken={} 完成登录
//
// 如果用户已登录，则直接重定向到 /me
func (r Controller) WeChatLogin(c *app.Context) result.Result {
	if c.Query("code") == "" {
		authPage, err := r.OfficialAccount.GetOauth().GetRedirectURL("http://gongfu.grianchan.com/wechat-login", "snsapi_userinfo", "login")
		if err != nil {
			return result.Err(fmt.Errorf("get redirect url: %w", err))
		}
		c.Redirect(http.StatusFound, authPage)
		return result.Ok(nil)
	} else {
		code := c.Query("code")
		accessToken, err := r.OfficialAccount.GetOauth().GetUserAccessToken(code)
		if err != nil {
			return result.Err(fmt.Errorf("get oauth: %w", err))
		}
		uInfo, err := r.OfficialAccount.GetOauth().GetUserInfo(accessToken.AccessToken, accessToken.OpenID, "zh_CN")
		if err != nil {
			return result.Err(fmt.Errorf("get oauth user info: %w", err))
		}
		user, err := r.Store.GetUserByOpenID(context.TODO(), uInfo.OpenID)
		if err != nil {
			return result.Err(fmt.Errorf("get user by open id: %w", err))
		}
		if user == nil {
			user = &model.User{
				OpenID:     &uInfo.OpenID,
				UnionID:    &uInfo.Unionid,
				Nickname:   uInfo.Nickname,
				Sex:        uInfo.Sex,
				Province:   uInfo.Province,
				City:       uInfo.City,
				Country:    uInfo.Country,
				HeadImgURL: uInfo.HeadImgURL,
				UUID:       util.UUID(),
			}
			if err := r.Store.CreateUser(context.TODO(), user); err != nil {
				return result.Err(fmt.Errorf("create user: %w", err))
			}
		}
		jwtToken, err := r.Token.GetAccessToken(user.ID)
		if err != nil {
			return result.Err(err)
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?accessToken=%s", wechatLoginURL, jwtToken))
		return result.Ok(nil)
	}
}

type MeResponse struct {
	ID         uint    `json:"id"`
	OpenID     *string `json:"openid"`
	Phone      *string `json:"phone"`
	Nickname   string  `json:"nickname"`
	HeadImgURL string  `json:"headimgurl"`
	Role       string  `json:"role"`
	UUID       string  `json:"uuid"`
}

func (r Controller) Me(c *app.Context) result.Result {
	var err error
	avatarURL := c.User.HeadImgURL
	if !strings.HasPrefix(avatarURL, "http") {
		avatarURL = r.Storage.GetPublicVisitURL(context.Background(), avatarURL) + "!avatar"
		if err != nil {
			return result.Err(err)
		}
	}
	return result.Ok(MeResponse{
		ID:         c.User.ID,
		OpenID:     c.User.OpenID,
		Phone:      c.User.Phone,
		Nickname:   c.User.Nickname,
		HeadImgURL: avatarURL,
		Role:       c.User.Role,
		UUID:       c.User.UUID,
	})
}

// EditMe 编辑用户信息
func (r Controller) EditMe(c *app.Context) result.Result {
	req := struct {
		AvatarKey string `json:"avatarKey" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		return result.Err(err)
	}
	// 检查储存中文件是否已经上传
	if exists, err := r.Storage.KeyExists(context.TODO(), req.AvatarKey); err != nil {
		return result.Err(err)
	} else if !exists {
		c.Message(http.StatusBadRequest, fmt.Sprintf("key %s no exists", req.AvatarKey))
		return result.Err(nil)
	}
	c.User.HeadImgURL = req.AvatarKey
	if err := r.Store.UpdateUser(context.TODO(), &c.User); err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r Controller) Login(c *app.Context) result.Result {
	req := struct {
		Phone string `binding:"required"` // 手机号
		Code  string `binding:"required"` // 验证码
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Message(http.StatusBadRequest, fmt.Sprintf("bad request params: %s", err))
		return result.Err(nil)
	}
	user, err := r.Store.GetUserByPhone(context.TODO(), req.Phone)
	if err != nil {
		return result.Err(err)
	}
	if user == nil {
		c.Message(http.StatusBadRequest, "user not exist")
		return result.Err(nil)
	}
	// 检查验证码
	if !r.Config.IsProd() && req.Code == "2222" {
		// 通过
		jwtToken, err := r.Token.GetAccessToken(user.ID)
		if err != nil {
			return result.Err(err)
		}
		return result.Ok(gin.H{
			"accessToken": jwtToken,
		})
	} else {
		// 验证码错误
		c.Message(http.StatusBadRequest, "code error")
		return result.Err(nil)
	}
}
