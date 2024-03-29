// Package date check in 相关的时间计算
// today 标准为凌晨5点至次日凌晨5点
package date

import "time"

// 一天的开始时间从凌晨 5 点开始
const dayStartHours = 5

// GetTodayDate 获取今天的日期 20210921
func GetTodayDate() string {
	return GetDateFromTime(time.Now())
}

func GetDayStartEnd(t time.Time) (time.Time, time.Time) {
	if t.Hour() >= dayStartHours {
		// 今天
		return time.Date(t.Year(), t.Month(), t.Day(), dayStartHours, 0, 0, 0, time.Local),
			time.Date(t.Year(), t.Month(), t.Day()+1, dayStartHours, 0, 0, 0, time.Local)
	} else {
		// 昨天
		return time.Date(t.Year(), t.Month(), t.Day()-1, dayStartHours, 0, 0, 0, time.Local),
			time.Date(t.Year(), t.Month(), t.Day(), dayStartHours, 0, 0, 0, time.Local)
	}
}

// GetTodayStartEnd 获取今天的开始时间结束时间
func GetTodayStartEnd() (time.Time, time.Time) {
	return GetDayStartEnd(time.Now())
}

// GetDateFromTime 从时间中获取日期
func GetDateFromTime(t time.Time) string {
	if t.Hour() >= dayStartHours {
		// 今天
		return t.Format("20060102")
	} else {
		// 昨天
		return t.AddDate(0, 0, -1).Format("20060102")
	}
}
