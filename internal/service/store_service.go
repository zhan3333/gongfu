package service

import (
	"context"
	"fmt"
	"gongfu/internal/model"
	"gongfu/pkg/date"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)
import _ "gorm.io/driver/mysql"

type Store interface {
	// CheckIn 打卡相关查询
	CheckIn
	// CheckInCount 打卡统计相关查询
	CheckInCount
	// User 用户相关查询
	User
}

type CheckIn interface {
	CreateCheckIn(ctx context.Context, checkIn *model.CheckIn) error
	GetCheckIn(ctx context.Context, options GetCheckInOption) (*model.CheckIn, error)
	// GetCheckInRankNum 获取 checkIn 的排名
	GetCheckInRankNum(ctx context.Context, checkIn *model.CheckIn) (int64, error)
}

type User interface {
	UserIDExists(ctx context.Context, userID uint) (bool, error)
	OpenIDExists(ctx context.Context, openID string) (bool, error)
	GetUserByOpenID(ctx context.Context, openID string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUser(ctx context.Context, userID uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUsersMap(ctx context.Context, userIDs []uint) (UsersMap, error)
	GetCoach(ctx context.Context, userID uint) (*model.Coach, error)
}

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

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func (s DBStore) GetCoach(ctx context.Context, userID uint) (*model.Coach, error) {
	coach := model.Coach{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&coach).Error
	if err != nil {
		return nil, err
	}
	if coach.ID == 0 {
		return nil, nil
	}
	return &coach, nil
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
	err := s.DB.WithContext(ctx).Where("check_in_continuous != ?", 0).Order("`check_in_continuous` desc").Limit(int(length)).Find(&counts).Error
	return counts, err
}

type GetCheckInHistoriesOptions struct {
	UserID    uint
	Length    uint
	Unique    bool
	StartDate string
	EndDate   string
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

func (s DBStore) GetUsersMap(ctx context.Context, userIDs []uint) (UsersMap, error) {
	var users = []*model.User{}
	err := s.DB.WithContext(ctx).Where("id in (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	ret := map[uint]*model.User{}
	for _, user := range users {
		ret[user.ID] = user
	}
	return ret, nil
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

func (s DBStore) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("phone = ?", phone).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) UpdateUser(ctx context.Context, user *model.User) error {
	return s.DB.WithContext(ctx).Save(user).Error
}

func (s DBStore) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("id = ?", userID).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) GetUserByOpenID(ctx context.Context, openID string) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("open_id = ?", openID).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func NewDBStore(DB *gorm.DB) Store {
	return &DBStore{DB: DB}
}

func (s DBStore) UserIDExists(ctx context.Context, userID uint) (bool, error) {
	count := int64(0)
	err := s.DB.WithContext(ctx).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s DBStore) OpenIDExists(ctx context.Context, openID string) (bool, error) {
	count := int64(0)
	err := s.DB.WithContext(ctx).Where("open_id = ?", openID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s DBStore) CreateUser(ctx context.Context, user *model.User) error {
	return s.DB.WithContext(ctx).Create(&user).Error
}
