// Package time 通用时间处理函数
package time

import "time"

// 时间布局.
const (
	DefaultLayout string = "2006-01-02 15:04:05" // 默认布局,秒布局
	DayLayout     string = "2006-01-02"          // 日期布局
)

// Parse 时间字符串转->时间结构体, layout 空采用默认布局
func Parse(value, layout string) (time.Time, error) {
	if layout == "" {
		layout = DefaultLayout
	}
	return time.ParseInLocation(layout, value, time.Local)
}

// GetTimeZero 获取当天零点时间
func GetTimeZero(t time.Time) time.Time {
	return time.Date(t.Local().Year(), t.Local().Month(), t.Local().Day(), 0, 0, 0, 0, time.Local)
}

// GetDiffDay 获取两个时间相差的天数, 0 为同一天，正数 t1 > t2, 负数 t1 < t2
func GetDiffDay(t1, t2 time.Time) int {
	return int(GetTimeZero(t1).Sub(GetTimeZero(t2)).Hours() / 24)
}

// GetDiffDayByTs 获取两个时间相差的天数，根据时间戳计算
func GetDiffDayByTs(t1, t2 int64) int {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)
	return GetDiffDay(time1, time2)
}

// GetMonthFirstDayZero 获取该时间月份第一天零点
func GetMonthFirstDayZero(t time.Time) time.Time {
	return GetTimeZero(t.AddDate(0, 0, -t.Day()+1))
}

// GetMondayOfThisWeek 取本周一零点零时
func GetMondayOfThisWeek(t time.Time) time.Time {
	// Sunday = 0, ... , Saturday = 6
	offset := int(time.Monday - t.Weekday())
	// 本周周一是第一天，要保证所有offset是负的
	if offset > 0 {
		offset -= 7
	}
	return GetTimeZero(t).AddDate(0, 0, offset)
}

// IsWeekend 是否是周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Sunday || weekday == time.Saturday
}

// IsSameWeek 判断是否一个周
func IsSameWeek(t1, t2 time.Time) bool {
	m1 := GetMondayOfThisWeek(t1)
	m2 := GetMondayOfThisWeek(t2)
	return m1.Equal(m2)
}

// IsSameDay 判断两个时间是否同一天
func IsSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
