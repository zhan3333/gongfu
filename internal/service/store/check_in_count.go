package store

import (
	"context"
	"fmt"
	"gongfu/internal/model"
	"gongfu/pkg/date"
	"time"
)

type CheckInCount interface {
	// GetCheckInCountTop 获取打卡次数排行榜
	GetCheckInCountTop(ctx context.Context, length uint) ([]*model.CheckInCount, error)
	// GetCheckInContinuousTop 获取连续打卡次数排行榜
	GetCheckInContinuousTop(ctx context.Context, length uint) ([]*model.CheckInCount, error)
	// GetCheckInHistories 获取指定用户打卡历史
	// unique: 是否只取每天最后一次打卡
	GetCheckInHistories(ctx context.Context, options GetCheckInHistoriesOptions) ([]*model.CheckIn, error)
	// GetCheckInTop 获取指定日期下的打卡排行 (打卡时间排名)
	GetCheckInTop(ctx context.Context, date string) ([]*model.CheckIn, error)
}

type GetCheckInHistoriesOptions struct {
	UserID    uint
	Length    uint
	Unique    bool
	StartDate string
	EndDate   string
}

// UpdateCheckInCount 更新统计信息
// 默认为今天未打卡，进行更新
func (s DBStore) UpdateCheckInCount(ctx context.Context, userID uint, isYesterdayCheckIn bool) error {
	count := &model.CheckInCount{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Find(count).Error
	if err != nil {
		return fmt.Errorf("query check in count: %w", err)
	}
	if count.ID == 0 {
		count = &model.CheckInCount{
			UserID:            userID,
			CheckInCount:      0,
			CheckInContinuous: 0,
		}
	}
	count.CheckInCount++
	if isYesterdayCheckIn {
		count.CheckInContinuous++
	} else {
		count.CheckInContinuous = 1
	}
	if err := s.DB.WithContext(ctx).Save(count).Error; err != nil {
		return fmt.Errorf("save check in count: %w", err)
	}
	return nil
}

func (s DBStore) GetCheckInTop(ctx context.Context, date string) ([]*model.CheckIn, error) {
	var checkIns = []*model.CheckIn{}
	err := s.DB.WithContext(ctx).Where(
		"id in (?)",
		s.DB.Model(&model.CheckInDay{}).Select("check_in_id").Where("date = ?", date),
	).Order("created_at").Find(&checkIns).Error
	if err != nil {
		return nil, err
	}
	if len(checkIns) == 0 {
		return []*model.CheckIn{}, nil
	}
	return checkIns, nil
}

func (s DBStore) GetCheckInCountTop(ctx context.Context, length uint) ([]*model.CheckInCount, error) {
	var counts = []*model.CheckInCount{}
	err := s.DB.WithContext(ctx).
		Where("check_in_count != ?", 0).
		Order("`check_in_count` desc").
		Limit(int(length)).
		Find(&counts).Error
	return counts, err
}

func (s DBStore) GetCheckInContinuousTop(ctx context.Context, length uint) ([]*model.CheckInCount, error) {
	var counts = []*model.CheckInCount{}
	// 更新时间需要在今日
	start, end := date.GetTodayStartEnd()
	err := s.DB.WithContext(ctx).
		Where("check_in_continuous != ?", 0).
		Where("updated_at between ? and ?", start, end).
		Order("`check_in_continuous` desc").
		Limit(int(length)).Find(&counts).Error
	return counts, err
}

func (s DBStore) GetCheckInHistories(ctx context.Context, options GetCheckInHistoriesOptions) ([]*model.CheckIn, error) {
	var checkIns = []*model.CheckIn{}
	if options.Unique {
		subQuery := s.DB.Model(&model.CheckInDay{}).Select("check_in_id").Where("user_id = ?", options.UserID)
		if options.StartDate != "" && options.EndDate != "" {
			subQuery = subQuery.Where("`date` between ? and ?", options.StartDate, options.EndDate)
		}
		query := s.DB.WithContext(ctx).Where("id in (?)", subQuery).
			Order("created_at desc")
		if options.Length != 0 {
			query = query.Limit(int(options.Length))
		}
		err := query.Find(&checkIns).Error
		return checkIns, err
	} else {
		query := s.DB.WithContext(ctx).Where("user_id = ?", options.UserID)
		if options.StartDate != "" && options.EndDate != "" {
			start, err := time.Parse("20060102", options.StartDate)
			if err != nil {
				return nil, fmt.Errorf("parse start date %s: %w", options.StartDate, err)
			}
			end, err := time.Parse("20060102", options.EndDate)
			if err != nil {
				return nil, fmt.Errorf("parse start date %s: %w", options.EndDate, err)
			}
			start = start.Add(5 * time.Hour).Add(-1 * time.Second)
			end = end.AddDate(0, 0, 0).Add(5 * time.Hour)
			query.Where("created_at between ? and ?", start, end)
		}
		if options.Length != 0 {
			query = query.Limit(int(options.Length))
		}
		err := query.Order("created_at desc").Find(&checkIns).Error
		return checkIns, err
	}
}
