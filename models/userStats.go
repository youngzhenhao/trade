package models

type UserStats struct {
	QueryTime            string    `json:"查询时间" yaml:"查询时间"`
	TotalUser            uint64    `json:"用户总数" yaml:"用户总数"`
	NewUserTodayNum      uint64    `json:"今日新用户数" yaml:"今日新用户数"`
	DailyActiveUserNum   uint64    `json:"日活跃用户数(DAU)" yaml:"日活跃用户数"`
	MonthlyActiveUserNum uint64    `json:"月活跃用户数(MAU)" yaml:"月活跃用户数"`
	NewUserToday         *[]User   `json:"今日新用户,omitempty" yaml:"今日新用户"`
	DailyActiveUser      *[]User   `json:"日活跃用户,omitempty" yaml:"日活跃用户"`
	MonthlyActiveUser    *[]User   `json:"月活跃用户,omitempty" yaml:"月活跃用户"`
	ErrorInfos           *[]string `json:"错误信息" yaml:"错误信息"`
}
