package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/model"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
	"net/http"
)

type Middlewares interface {
	Auth() gin.HandlerFunc
	Role(roleName string) gin.HandlerFunc
}

type middlewares struct {
	Token service.Token
	Store store.Store
}

func NewMiddlewares(token service.Token, store store.Store) Middlewares {
	return &middlewares{Token: token, Store: store}
}

func (m middlewares) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "no auth header"})
			return
		}

		if userID, err := m.Token.ParseAccessToken(token); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": err.Error()})
			return
		} else {
			c.Set("user_id", userID)
			if user, err := m.Store.GetUser(context.TODO(), userID); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("find user error: %s", err.Error())})
				return
			} else {
				c.Set("user", *user)
			}
		}
		c.Next()
	}
}

// Role 验证登录用户是否属于某个角色
func (m middlewares) Role(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, ok := c.Get("user"); ok {
			if user, ok := user.(model.User); ok {
				if user.HasRole(roleName) {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "role must contain " + roleName})
		return
	}
}
