package controller

import (
	"context"
	"fmt"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"net/http"
)

var meURL = "/web/me"

// Login 公众号登录
// /login 时会 302 到微信授权登录页面, 用户确认授权后返回到 /login?code={} 地址
// /login?code={} 时处理 code 为 accessToken 获取微信用户信息
// 获取微信用户信息后，根据 openID 查找用户并生成 JWT token
// 生成 JWT token 后 302 到 /me?accessToken={} 完成登录
//
// 如果用户已登录，则直接重定向到 /me
func (r Controller) Login(c *app.Context) result.Result {
	if c.Query("code") == "" {
		authPage, err := r.OfficialAccount.GetOauth().GetRedirectURL("http://gongfu.grianchan.com/login", "snsapi_userinfo", "login")
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
			}
			if err := r.Store.CreateUser(context.TODO(), user); err != nil {
				return result.Err(fmt.Errorf("create user: %w", err))
			}
		}
		jwtToken, err := r.Token.GetAccessToken(user.ID)
		if err != nil {
			return result.Err(err)
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("%s?accessToken=%s", meURL, jwtToken))
		return result.Ok(nil)
	}
}

func (r Controller) Me(c *app.Context) result.Result {
	return result.Ok(c.User)
}
