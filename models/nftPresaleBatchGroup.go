package models

import (
	"gorm.io/gorm"
	"time"
)

type NftPresaleBatchGroup struct {
	gorm.Model
	GroupKey     string `json:"group_key" gorm:"type:varchar(255);index"`
	GroupName    string `json:"group_name" gorm:"type:varchar(255);index"`
	SoldNumber   int    `json:"sold_number"`
	Supply       int    `json:"supply"`
	LowestPrice  int    `json:"lowest_price"`
	HighestPrice int    `json:"highest_price"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	Info         string `json:"info"`
}

type NftPresaleBatchGroupSetRequest struct {
	GroupKey string `json:"group_key" gorm:"type:varchar(255);index"`
	// not been used
	Supply    uint   `json:"supply"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	Info      string `json:"info"`
}

type NftPresaleBatchGroupSimplified struct {
	ID           uint `gorm:"primarykey"`
	UpdatedAt    time.Time
	GroupKey     string `json:"group_key" gorm:"type:varchar(255);index"`
	GroupName    string `json:"group_name" gorm:"type:varchar(255);index"`
	SoldNumber   int    `json:"sold_number"`
	Supply       int    `json:"supply"`
	LowestPrice  int    `json:"lowest_price"`
	HighestPrice int    `json:"highest_price"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	Info         string `json:"info"`
}

type NftPresaleBatchGroupLaunchRequest struct {
	BatchGroupSetRequest  NftPresaleBatchGroupSetRequest
	NftPresaleSetRequests *[]NftPresaleSetRequest
}

type NftPresaleBatchGroupState int

const (
	NftPresaleBatchGroupStateAll NftPresaleBatchGroupState = iota
	NftPresaleBatchGroupStateSelling
	NftPresaleBatchGroupStateNotStart
	NftPresaleBatchGroupStateEnd
)
