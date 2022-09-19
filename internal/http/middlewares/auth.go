package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/service"
	"net/http"
)

type Middlewares interface {
	Auth() gin.HandlerFunc
}

type middlewares struct {
	Token service.Token
	Store service.Store
}

func NewMiddlewares(token service.Token, store service.Store) Middlewares {
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
