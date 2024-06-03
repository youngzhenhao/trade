package models

import "gorm.io/gorm"

type ScheduledTask struct {
	gorm.Model
	Name           string `gorm:"type:varchar(255);not null" json:"name"`
	CronExpression string `gorm:"type:varchar(100);not null" json:"cron_expression"`
	FunctionName   string `gorm:"type:varchar(100);not null" json:"function_name"`
	Package        string `gorm:"type:varchar(255);not null" json:"package"`
	Status         int16  `gorm:"default:1" json:"status"`
}
