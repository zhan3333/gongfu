package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/internal/service/store"
	"gongfu/pkg/date"
	"net/http"
	"time"
)

// PostCheckIn 打卡
func (r UseCase) PostCheckIn(c *app.Context) result.Result {
	req := struct {
		Key      string `json:"key" binding:"required"`
		FileName string `json:"fileName" binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		return result.Err(err)
	}
	// 检查储存中文件是否已经上传
	if exists, err := r.Storage.KeyExists(context.TODO(), req.Key); err != nil {
		return result.Err(err)
	} else if !exists {
		c.Message(http.StatusBadRequest, fmt.Sprintf("key %s no exists", req.Key))
		return result.Err(nil)
	}
	err := r.Store.CreateCheckIn(context.TODO(), &model.CheckIn{
		Key:       req.Key,
		UserID:    c.UserID,
		FileName:  req.FileName,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

type CheckInResp struct {
	ID         uint   `json:"id"`
	CreatedAt  int64  `json:"createdAt"`
	Url        string `json:"url"`
	Key        string `json:"key"`
	UserName   string `json:"userName"`
	UserID     uint   `json:"userID"`
	HeadImgUrl string `json:"headImgUrl"`
	Date       string `json:"date"`
	DayRank    int64  `json:"dayRank"`
	UserUUID   string `json:"userUUID"`
}

// GetTodayCheckIn 获取当前用户的今日打卡记录
func (r UseCase) GetTodayCheckIn(c *app.Context) result.Result {
	checkIn, err := r.Store.GetCheckIn(context.TODO(), store.GetCheckInOption{UserID: c.UserID, Date: date.GetTodayDate()})
	if err != nil {
		return result.Err(err)
	}
	if checkIn == nil {
		return result.Ok(gin.H{
			"exists": false,
		})
	}
	user, err := r.Store.GetUser(context.TODO(), checkIn.UserID)
	if err != nil {
		return result.Err(fmt.Errorf("get user: %w", err))
	}
	var userName, headImgUrl, userID = "", "", uint(0)
	if user != nil {
		userName, headImgUrl, userID = user.Nickname, r.Storage.GetHeadImageVisitURL(user.HeadImgURL), user.ID
	}

	visitUrl, err := r.Storage.GetVisitURL(context.TODO(), checkIn.Key)
	if err != nil {
		return result.Err(fmt.Errorf("get visit url: %w", err))
	}

	// 日排名
	dayRank, err := r.Store.GetCheckInRankNum(context.TODO(), checkIn)
	if err != nil {
		return result.Err(fmt.Errorf("get rank num: %w", err))
	}
	return result.Ok(gin.H{
		"exists": true,
		"checkIn": CheckInResp{
			ID:         checkIn.ID,
			CreatedAt:  checkIn.CreatedAt.Unix(),
			Url:        visitUrl,
			Key:        checkIn.Key,
			UserName:   userName,
			UserID:     userID,
			HeadImgUrl: headImgUrl,
			Date:       date.GetDateFromTime(checkIn.CreatedAt),
			DayRank:    dayRank,
		},
	})
}

// GetCheckIn 根据 key 获取打卡详情
func (r UseCase) GetCheckIn(c *app.Context) result.Result {
	checkInKey := c.Param("key")
	if checkInKey == "" {
		c.Message(http.StatusBadRequest, "invalid check in key: "+checkInKey)
		return result.Err(nil)
	}
	checkIn, err := r.Store.GetCheckIn(context.TODO(), store.GetCheckInOption{Key: checkInKey})
	if err != nil {
		return result.Err(fmt.Errorf("get check in: %w", err))
	}
	if checkIn == nil {
		c.Message(http.StatusNotFound, "check in not found")
		return result.Err(nil)
	}
	user, err := r.Store.GetUser(context.TODO(), checkIn.UserID)
	if err != nil {
		return result.Err(fmt.Errorf("get user: %w", err))
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "user not found"})
		return result.Err(nil)
	}
	// 日排名
	dayRank, err := r.Store.GetCheckInRankNum(context.TODO(), checkIn)
	if err != nil {
		return result.Err(fmt.Errorf("get rank num: %w", err))
	}
	visitUrl, err := r.Storage.GetVisitURL(context.TODO(), checkIn.Key)
	if err != nil {
		return result.Err(fmt.Errorf("get visit url: %w", err))
	}
	return result.Ok(CheckInResp{
		ID:         checkIn.ID,
		CreatedAt:  checkIn.CreatedAt.Unix(),
		Url:        visitUrl,
		Key:        checkIn.Key,
		UserName:   user.Nickname,
		UserID:     user.ID,
		HeadImgUrl: r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
		DayRank:    dayRank,
		UserUUID:   user.UUID,
	})
}

// GetCheckInTop 获取打卡排行
func (r UseCase) GetCheckInTop(c *app.Context) result.Result {
	req := struct {
		Date string `json:"date" form:"date"`
	}{}
	var resp []struct {
		CreatedAt  int64  `json:"createdAt"`
		ID         uint   `json:"id"`
		UserName   string `json:"userName"`
		HeadImgUrl string `json:"headImgUrl"`
		Key        string `json:"key"`
		UserUUID   string `json:"userUUID"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		return result.Err(err)
	}
	if req.Date == "" {
		req.Date = date.GetTodayDate()
	}
	if checkInList, err := r.Store.GetCheckInTop(context.TODO(), req.Date); err != nil {
		return result.Err(err)
	} else {
		var userIds []uint
		for _, checkIn := range checkInList {
			userIds = append(userIds, checkIn.UserID)
		}
		usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
		if err != nil {
			return result.Err(fmt.Errorf("get users map: %w", err))
		}
		for _, checkIn := range checkInList {
			user := usersMap.DefaultGet(checkIn.UserID)
			resp = append(resp, struct {
				CreatedAt  int64  `json:"createdAt"`
				ID         uint   `json:"id"`
				UserName   string `json:"userName"`
				HeadImgUrl string `json:"headImgUrl"`
				Key        string `json:"key"`
				UserUUID   string `json:"userUUID"`
			}{
				CreatedAt:  checkIn.CreatedAt.Unix(),
				ID:         checkIn.ID,
				UserName:   user.Nickname,
				HeadImgUrl: r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
				Key:        checkIn.Key,
				UserUUID:   user.UUID,
			})
		}
		return result.Ok(resp)
	}
}

type CheckInCountItem struct {
	ID         uint   `json:"id"`
	UserName   string `json:"userName"`
	UserID     uint   `json:"userID"`
	HeadImgUrl string `json:"headImgUrl"`
	// 打卡总次数
	CheckInCount uint   `json:"checkInCount"`
	UserUUID     string `json:"userUUID"`
}

// GetCheckInCountTop 获取总打卡次数排行榜
func (r UseCase) GetCheckInCountTop(_ *app.Context) result.Result {
	checkInList, err := r.Store.GetCheckInCountTop(context.TODO(), 10)
	if err != nil {
		return result.Err(fmt.Errorf("get check in count top: %w", err))
	}
	var userIDs []uint
	for _, item := range checkInList {
		userIDs = append(userIDs, item.UserID)
	}
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIDs)
	if err != nil {
		return result.Err(fmt.Errorf("get users map: %w", err))
	}
	var res []CheckInCountItem
	for _, item := range checkInList {
		user := usersMap.DefaultGet(item.UserID)
		res = append(res, CheckInCountItem{
			ID:           item.ID,
			UserName:     user.Nickname,
			UserID:       item.UserID,
			HeadImgUrl:   r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
			CheckInCount: item.CheckInCount,
			UserUUID:     user.UUID,
		})
	}
	return result.Ok(res)
}

type CheckInContinuousItem struct {
	ID         uint   `json:"id"`
	UserName   string `json:"userName"`
	UserID     uint   `json:"userID"`
	HeadImgUrl string `json:"headImgUrl"`
	// 连续打卡次数
	CheckInContinuous uint   `json:"checkInContinuous"`
	UserUUID          string `json:"userUUID"`
}

// GetCheckInContinuousTop 获取连续打卡排行榜
func (r UseCase) GetCheckInContinuousTop(_ *app.Context) result.Result {
	checkInList, err := r.Store.GetCheckInContinuousTop(context.TODO(), 10)
	if err != nil {
		return result.Err(fmt.Errorf("get check in count top: %w", err))
	}
	var userIDs []uint
	for _, item := range checkInList {
		userIDs = append(userIDs, item.UserID)
	}
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIDs)
	if err != nil {
		return result.Err(fmt.Errorf("get users map: %w", err))
	}
	var res []CheckInContinuousItem
	for _, item := range checkInList {
		user := usersMap.DefaultGet(item.UserID)
		res = append(res, CheckInContinuousItem{
			ID:                item.ID,
			UserName:          user.Nickname,
			UserID:            item.UserID,
			HeadImgUrl:        r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
			CheckInContinuous: item.CheckInContinuous,
			UserUUID:          user.UUID,
		})
	}
	return result.Ok(res)
}

// GetCheckInHistories 获取打卡历史
func (r UseCase) GetCheckInHistories(c *app.Context) result.Result {
	req := struct {
		StartDate string `json:"startDate" form:"startDate" binding:"required"`
		EndDate   string `json:"endDate" form:"endDate" binding:"required"`
		UserID    uint   `json:"userID" form:"userID" binding:"required"`
		// 是否根据 date 去重
		Unique *bool `json:"unique" form:"unique" binding:"omitempty"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Message(http.StatusBadRequest, fmt.Sprintf("invalid request"))
		return result.Err(nil)
	}

	unique := true
	if req.Unique != nil {
		unique = *req.Unique
	}
	length := uint(0)
	if req.StartDate == "" || req.EndDate == "" {
		length = 10
	}

	items, err := r.Store.GetCheckInHistories(context.TODO(), store.GetCheckInHistoriesOptions{
		UserID:    req.UserID,
		Length:    length,
		Unique:    unique,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	if err != nil {
		return result.Err(fmt.Errorf("get historires: %w", err))
	}
	var userIDs []uint
	for _, item := range items {
		userIDs = append(userIDs, item.UserID)
	}
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIDs)
	if err != nil {
		return result.Err(fmt.Errorf("get users map: %w", err))
	}
	var res []CheckInResp
	for _, item := range items {
		visitUrl, err := r.Storage.GetVisitURL(context.TODO(), item.Key)
		if err != nil {
			return result.Err(fmt.Errorf("get %s visit url: %w", item.Key, err))
		}
		res = append(res, CheckInResp{
			ID:         item.ID,
			CreatedAt:  item.CreatedAt.Unix(),
			Url:        visitUrl,
			Key:        item.Key,
			UserName:   usersMap.DefaultGet(item.UserID).Nickname,
			UserID:     item.UserID,
			HeadImgUrl: r.Storage.GetHeadImageVisitURL(usersMap.DefaultGet(item.UserID).HeadImgURL),
			Date:       date.GetDateFromTime(item.CreatedAt),
		})
	}
	if len(res) == 0 {
		return result.Ok([]CheckInResp{})
	}
	return result.Ok(res)
}
