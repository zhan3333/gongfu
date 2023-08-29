package store

import (
	"context"
	"fmt"
	"gongfu/internal/model"
	"gongfu/pkg/date"
)

type CheckIn interface {
	CreateCheckIn(ctx context.Context, checkIn *model.CheckIn) error
	GetCheckIn(ctx context.Context, options GetCheckInOption) (*model.CheckIn, error)
	// GetCheckInRankNum 获取 checkIn 的排名
	GetCheckInRankNum(ctx context.Context, checkIn *model.CheckIn) (int64, error)
}

type GetCheckInOption struct {
	CheckInID uint
	UserID    uint
	Date      string
	Key       string
}

func (s DBStore) GetCheckIn(ctx context.Context, options GetCheckInOption) (*model.CheckIn, error) {
	if options.CheckInID == 0 && options.Date == "" && options.Key == "" && options.UserID == 0 {
		return nil, fmt.Errorf("options not allowed empty")
	}
	var checkIn = &model.CheckIn{}
	childQuery := s.DB.WithContext(ctx).Model(&model.CheckInDay{}).Select("check_in_id")
	if options.CheckInID != 0 {
		childQuery = childQuery.Where("check_in_id = ?", options.CheckInID)
	}
	if options.UserID != 0 {
		childQuery = childQuery.Where("user_id = ?", options.UserID)
	}
	if options.Date != "" {
		childQuery = childQuery.Where("`date` = ?", options.Date)
	}
	query := s.DB.WithContext(ctx)
	if options.Key != "" {
		query = query.Where("`key` = ?", options.Key)
	}
	if options.CheckInID != 0 || options.UserID != 0 || options.Date != "" {
		query = query.Where("id in (?)", childQuery)
	}
	err := query.Find(&checkIn).Error
	if err != nil {
		return nil, err
	}
	if checkIn.ID == 0 {
		return nil, nil
	}
	return checkIn, nil
}

// CreateCheckIn 创建打卡记录
// 先检查今日是否已经打卡，未打卡则会增加 count
// 如果今日未打卡：昨日已打卡则 continuous + 1，昨日未打卡则置为 1
func (s DBStore) CreateCheckIn(ctx context.Context, checkIn *model.CheckIn) error {
	today := date.GetDateFromTime(checkIn.CreatedAt)
	yesterday := date.GetDateFromTime(checkIn.CreatedAt.AddDate(0, 0, -1))

	todayCheckIn, err := s.GetCheckIn(ctx, GetCheckInOption{UserID: checkIn.UserID, Date: today})
	if err != nil {
		return fmt.Errorf("query today check in: %w", err)
	}
	// 创建打卡记录
	if err := s.DB.WithContext(ctx).Create(checkIn).Error; err != nil {
		return fmt.Errorf("create check in: %w", err)
	}
	if err := s.UpdateCheckInDay(ctx, checkIn.UserID, today, checkIn.ID); err != nil {
		return fmt.Errorf("update check in day: %w", err)
	}
	// 更新统计
	if todayCheckIn == nil {
		yesterdayCheckIn, err := s.GetCheckIn(ctx, GetCheckInOption{UserID: checkIn.UserID, Date: yesterday})
		if err != nil {
			return fmt.Errorf("query yesterday check in: %w", err)
		}
		if err := s.UpdateCheckInCount(ctx, checkIn.UserID, yesterdayCheckIn != nil); err != nil {
			return fmt.Errorf("update check in count: %w", err)
		}
	}
	return nil
}

func (s DBStore) UpdateCheckInDay(ctx context.Context, userID uint, date string, checkInID uint) error {
	var day = &model.CheckInDay{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Where("`date` = ?", date).Find(day).Error
	if err != nil {
		return fmt.Errorf("query check in day: %w", err)
	}
	if day.ID == 0 {
		day = &model.CheckInDay{
			UserID:    userID,
			Date:      date,
			CheckInID: checkInID,
		}
	} else {
		day.CheckInID = checkInID
	}
	if err := s.DB.WithContext(ctx).Save(day).Error; err != nil {
		return fmt.Errorf("save check in day: %w", err)
	}
	return nil
}

func (s DBStore) GetCheckInRankNum(ctx context.Context, checkIn *model.CheckIn) (int64, error) {
	checkInDate := date.GetDateFromTime(checkIn.CreatedAt)
	checkInNum := int64(0)
	// 相同日期下，创建时间比传入 checkIn 早，且 userID 去重的 checkIn 数量
	q := s.DB.Model(&model.CheckInDay{}).
		Distinct("user_id").
		Where("date = ?", checkInDate).
		Where("created_at <= ?", checkIn.CreatedAt).
		Where("check_in_id != ?", checkIn.ID).
		Count(&checkInNum)
	if q.Error != nil {
		return 0, q.Error
	}
	return checkInNum + 1, nil
}
