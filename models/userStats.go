package models

type UserStats struct {
	QueryTime            string           `json:"查询时间" yaml:"查询时间"`
	TotalUser            uint64           `json:"用户总数" yaml:"用户总数"`
	SpecifiedDate        string           `json:"指定查询日期,omitempty" yaml:"指定查询日期"`
	NewUserTodayNum      uint64           `json:"今日新用户数" yaml:"今日新用户数"`
	DailyActiveUserNum   uint64           `json:"日活跃用户数(DAU)" yaml:"日活跃用户数"`
	MonthlyActiveUserNum uint64           `json:"月活跃用户数(MAU)" yaml:"月活跃用户数"`
	NewUserToday         *[]StatsUserInfo `json:"今日新用户,omitempty" yaml:"今日新用户"`
	DailyActiveUser      *[]StatsUserInfo `json:"日活跃用户,omitempty" yaml:"日活跃用户"`
	MonthlyActiveUser    *[]StatsUserInfo `json:"月活跃用户,omitempty" yaml:"月活跃用户"`
	ErrorInfos           *[]string        `json:"错误信息" yaml:"错误信息"`
}

type StatsUserInfo struct {
	ID                uint   `json:"用户ID" yaml:"用户ID"`
	CreatedAt         string `json:"用户创建时间" yaml:"用户创建时间"`
	UpdatedAt         string `json:"更新时间" yaml:"更新时间"`
	Username          string `json:"用户名;NpubKey;Nostr地址" yaml:"用户名;NpubKey;Nostr地址"`
	Status            string `json:"用户状态" yaml:"用户状态"`
	RecentIpAddresses string `json:"最近IP地址" yaml:"最近IP地址"`
	RecentLoginTime   string `json:"最近登录时间" yaml:"最近登录时间"`
}
