package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"net/http"
)

func (r Controller) CheckIn(c *app.Context) result.Result {
	req := struct {
		Key      string `json:"key" binding:"required"`
		FileName string `json:"fileName" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		return result.Err(err)
	}
	if exists, err := r.Storage.KeyExists(context.TODO(), req.Key); err != nil {
		return result.Err(err)
	} else if !exists {
		c.Message(http.StatusBadRequest, fmt.Sprintf("key %s no exists", req.Key))
		return result.Err(nil)
	}
	err := r.Store.CreateCheckIn(context.TODO(), &model.CheckIn{
		Key:      req.Key,
		UserID:   c.UserID,
		FileName: req.FileName,
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r Controller) GetTodayCheckIn(c *app.Context) result.Result {
	checkIn, err := r.Store.GetTodayCheckIn(context.TODO(), c.UserID)
	if err != nil {
		return result.Err(err)
	}
	if checkIn == nil {
		return result.Ok(gin.H{
			"exists": false,
		})
	}
	visitUrl, err := r.Storage.GetVisitURL(context.TODO(), checkIn.Key)
	if err != nil {
		return result.Err(fmt.Errorf("get visit url: %w", err))
	}
	return result.Ok(gin.H{
		"exists": true,
		"checkIn": gin.H{
			"id":        checkIn.ID,
			"createdAt": checkIn.CreatedAt.Unix(),
			"url":       visitUrl,
			"key":       checkIn.Key,
		},
	})
}

func (r Controller) GetCheckIn(c *app.Context) result.Result {
	checkInKey := c.Param("key")
	if checkInKey == "" {
		c.Message(http.StatusBadRequest, "invalid check in key: "+checkInKey)
		return result.Err(nil)
	}
	checkIn, err := r.Store.GetCheckInByKey(context.TODO(), checkInKey)
	if err != nil {
		return result.Err(fmt.Errorf("get check in: %w", err))
	}
	if checkIn == nil {
		c.Message(http.StatusNotFound, "check in not found")
		return result.Err(nil)
	}
	userName := " not found user"
	headImgUrl := ""
	user, err := r.Store.GetUser(context.TODO(), checkIn.UserID)
	if err != nil {
		return result.Err(fmt.Errorf("get user: %w", err))
	}
	if user != nil {
		headImgUrl = user.HeadImgURL
		userName = user.Nickname
	}
	visitUrl, err := r.Storage.GetVisitURL(context.TODO(), checkIn.Key)
	if err != nil {
		return result.Err(fmt.Errorf("get visit url: %w", err))
	}
	return result.Ok(gin.H{
		"id":         checkIn.ID,
		"createdAt":  checkIn.CreatedAt.Unix(),
		"url":        visitUrl,
		"fileName":   checkIn.FileName,
		"userName":   userName,
		"headImgUrl": headImgUrl,
	})
}

func (r Controller) GetCheckInTop(a *app.Context) result.Result {
	req := struct {
		StartAt int64 `form:"startAt" binding:"required"`
		EndAt   int64 `form:"endAt" binding:"required"`
	}{}
	var resp []struct {
		CreatedAt  int64  `json:"createdAt"`
		ID         uint   `json:"id"`
		UserName   string `json:"userName"`
		HeadImgUrl string `json:"headImgUrl"`
		Key        string `json:"key"`
	}
	if err := a.ShouldBindQuery(&req); err != nil {
		return result.Err(err)
	}
	if checkInList, err := r.Store.GetCheckInTop(context.TODO(), req.StartAt, req.EndAt); err != nil {
		return result.Err(err)
	} else {
		for _, checkIn := range checkInList {
			fmt.Println(checkIn)
		}
		var userIds []uint
		for _, checkIn := range checkInList {
			userIds = append(userIds, checkIn.UserID)
		}
		usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
		if err != nil {
			return result.Err(fmt.Errorf("get users map: %w", err))
		}
		for _, checkIn := range checkInList {
			userName := "not found user name"
			headImgUrl := ""
			if user, ok := usersMap[checkIn.UserID]; ok {
				userName = user.Nickname
				headImgUrl = user.HeadImgURL
			}
			resp = append(resp, struct {
				CreatedAt  int64  `json:"createdAt"`
				ID         uint   `json:"id"`
				UserName   string `json:"userName"`
				HeadImgUrl string `json:"headImgUrl"`
				Key        string `json:"key"`
			}{
				CreatedAt:  checkIn.CreatedAt.Unix(),
				ID:         checkIn.ID,
				UserName:   userName,
				HeadImgUrl: headImgUrl,
				Key:        checkIn.Key,
			})
		}
		return result.Ok(resp)
	}
}
