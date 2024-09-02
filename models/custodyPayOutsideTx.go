package models

import "gorm.io/gorm"

type PayOutsideTx struct {
	gorm.Model
	TxHash     string             `gorm:"column:tx_hash;index:idx_tx_hash;unique" json:"txHash"`
	Timestamp  int64              `gorm:"column:timestamp;" json:"timestamp"`
	HeightHint uint32             `gorm:"column:height_hint;" json:"heightHint"`
	ChainFees  int64              `gorm:"column:chain_fees;" json:"chainFees"`
	InputsNum  uint               `gorm:"column:inputs_num;" json:"inputsNum"`
	OutputsNum uint               `gorm:"column:outputs_num;" json:"outputsNum"`
	Status     PayOutsideTxStatus `gorm:"column:status;" json:"status"`
}

func (PayOutsideTx) TableName() string {
	return "user_out_inside_tx"
}

type PayOutsideTxStatus uint

const (
	PayOutsideStatusTXPending PayOutsideTxStatus = iota
	PayOutsideStatusTXSuccess
	PayOutsideStatusTXFailed
)
