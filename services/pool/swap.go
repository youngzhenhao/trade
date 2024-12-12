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

type SwapFeeType int64

const (
	SwapFee6ThousandsNotSat SwapFeeType = iota
	SwapFee6Thousands
	SwapFee20Sat
)

type PoolSwapRecord struct {
	gorm.Model
	PairId                   uint           `json:"pair_id" gorm:"index"`
	Username                 string         `json:"username" gorm:"type:varchar(255);index"`
	TokenIn                  string         `json:"token_in" gorm:"type:varchar(255);index"`
	TokenOut                 string         `json:"token_out" gorm:"type:varchar(255);index"`
	AmountIn                 string         `json:"amount_in" gorm:"type:varchar(255);index"`
	AmountOut                string         `json:"amount_out" gorm:"type:varchar(255);index"`
	ReserveIn                string         `json:"reserve_in" gorm:"type:varchar(255);index"`
	ReserveOut               string         `json:"reserve_out" gorm:"type:varchar(255);index"`
	TokenInTransferRecordId  uint           `json:"token_in_transfer_record_id" gorm:"index"`
	TokenOutTransferRecordId uint           `json:"token_out_transfer_record_id" gorm:"index"`
	SwapFee                  string         `json:"swap_fee" gorm:"type:varchar(255);index"`
	SwapFeeType              SwapFeeType    `json:"swap_fee_type" gorm:"index"`
	SwapRecordType           SwapRecordType `json:"swap_record_type" gorm:"index"`
}

func newSwapRecord(pairId uint, username string, tokenIn string, tokenOut string, amountIn string, amountOut string, reserveIn string, reserveOut string, tokenInTransferRecordId uint, tokenOutTransferRecordId uint, swapFee string, swapFeeType SwapFeeType, swapRecordType SwapRecordType) (swapRecord *PoolSwapRecord, err error) {
	if pairId <= 0 {
		return new(PoolSwapRecord), errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}
	return &PoolSwapRecord{
		PairId:                   pairId,
		Username:                 username,
		TokenIn:                  tokenIn,
		TokenOut:                 tokenOut,
		AmountIn:                 amountIn,
		AmountOut:                amountOut,
		ReserveIn:                reserveIn,
		ReserveOut:               reserveOut,
		TokenInTransferRecordId:  tokenInTransferRecordId,
		TokenOutTransferRecordId: tokenOutTransferRecordId,
		SwapFee:                  swapFee,
		SwapFeeType:              swapFeeType,
		SwapRecordType:           swapRecordType,
	}, nil
}

func createSwapRecord(tx *gorm.DB, pairId uint, username string, tokenIn string, tokenOut string, amountIn string, amountOut string, reserveIn string, reserveOut string, tokenInTransferRecordId uint, tokenOutTransferRecordId uint, swapFee string, swapFeeType SwapFeeType, swapRecordType SwapRecordType) (recordId uint, err error) {
	if pairId <= 0 {
		return 0, errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}
	var swapRecord *PoolSwapRecord
	swapRecord, err = newSwapRecord(pairId, username, tokenIn, tokenOut, amountIn, amountOut, reserveIn, reserveOut, tokenInTransferRecordId, tokenOutTransferRecordId, swapFee, swapFeeType, swapRecordType)
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "newSwapRecord")
	}
	err = tx.Model(&PoolSwapRecord{}).Create(&swapRecord).Error
	if err != nil {
		return 0, utils.AppendErrorInfo(err, "create swapRecord")
	}
	recordId = swapRecord.ID
	return recordId, nil
}

func calcSwapRecord(pairId uint, username string, tokenIn string, tokenOut string, amountIn string, amountOut string, reserveIn string, reserveOut string, swapFee string, swapFeeType SwapFeeType, swapRecordType SwapRecordType) (swapRecord *PoolSwapRecord, err error) {
	if pairId <= 0 {
		return new(PoolSwapRecord), errors.New("invalid pairId(" + strconv.FormatUint(uint64(pairId), 10) + ")")
	}
	swapRecord, err = newSwapRecord(pairId, username, tokenIn, tokenOut, amountIn, amountOut, reserveIn, reserveOut, 0, 0, swapFee, swapFeeType, swapRecordType)
	if err != nil {
		return new(PoolSwapRecord), utils.AppendErrorInfo(err, "newSwapRecord")
	}
	return swapRecord, nil
}
