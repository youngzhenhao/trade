package pool

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"trade/utils"
)

type SwapRecordType int64

const (
	SwapExactTokenNoPath SwapRecordType = iota
	SwapForExactTokenNoPath
)

type SwapRecord struct {
	gorm.Model
	PairId         uint           `json:"pair_id" gorm:"index"`
	Username       string         `json:"username" gorm:"type:varchar(255);index"`
	TokenIn        string         `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut       string         `json:"token_out" gorm:"type:varchar(255);index"`
	AmountIn       string         `json:"amount_in" gorm:"type:varchar(255);index"`
	AmountOut      string         `json:"amount_out" gorm:"type:varchar(255);index"`
	ReserveIn      string         `json:"reserve_in" gorm:"type:varchar(255);index"`
	ReserveOut     string         `json:"reserve_out" gorm:"type:varchar(255);index"`
	SwapRecordType SwapRecordType `json:"swap_record_type" gorm:"index"`
}

func NewSwapRecord(pairId uint, username string, tokenIn string, tokenOut string, amountIn string, amountOut string, reserveIn string, reserveOut string, swapRecordType SwapRecordType) (swapRecord *SwapRecord, err error) {
	if pairId <= 0 {
		return new(SwapRecord), errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}
	return &SwapRecord{
		PairId:         pairId,
		Username:       username,
		TokenIn:        tokenIn,
		TokenOut:       tokenOut,
		AmountIn:       amountIn,
		AmountOut:      amountOut,
		ReserveIn:      reserveIn,
		ReserveOut:     reserveOut,
		SwapRecordType: swapRecordType,
	}, nil
}

func CreateSwapRecord(tx *gorm.DB, pairId uint, username string, tokenIn string, tokenOut string, amountIn string, amountOut string, reserveIn string, reserveOut string, swapRecordType SwapRecordType) (err error) {
	if pairId <= 0 {
		return errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}
	var swapRecord *SwapRecord
	swapRecord, err = NewSwapRecord(pairId, username, tokenIn, tokenOut, amountIn, amountOut, reserveIn, reserveOut, swapRecordType)
	if err != nil {
		return utils.AppendErrorInfo(err, "NewSwapRecord")
	}
	err = tx.Model(&SwapRecord{}).Create(&swapRecord).Error
	if err != nil {
		return utils.AppendErrorInfo(err, "create swapRecord")
	}
	return nil
}
