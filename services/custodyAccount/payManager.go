package custodyAccount

import (
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/btc_channel"
)

// 托管账户划扣费用
func PayAmountToAdmin(payUserId uint, gasFee uint64) (uint, error) {
	e, err := btc_channel.NewBtcChannelEventByUserId(payUserId)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	id, err := btc_channel.PayFirLunchFee(e, gasFee)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return 0, err
	}
	return id, nil
}

func BackAmount(payInsideId uint) (uint, error) {
	missionId, err := CreateBackFeeMission(payInsideId)
	if err != nil {
		return 0, err
	}
	return missionId, nil
}

func CheckBackFeeMission(missionId uint) bool {
	return checkBackFeeMissionById(missionId)
}

// CheckPayInsideStatus 检查内部转账任务状态是否成功
func CheckPayInsideStatus(id uint) (bool, error) {
	p, err := btldb.ReadPayInside(id)
	if err != nil {
		return false, err
	}
	switch p.Status {
	case models.PayInsideStatusSuccess:
		return true, nil
	case models.PayInsideStatusFailed:
		return false, models.CustodyAccountPayInsideMissionFaild
	default:
		return false, models.CustodyAccountPayInsideMissionPending
	}
}

// IsAccountBalanceEnoughByUserId  判断账户余额是否足够
func IsAccountBalanceEnoughByUserId(userId uint, value uint64) bool {
	e, err := btc_channel.NewBtcChannelEventByUserId(userId)
	if err != nil {
		btlLog.CUST.Error("PayAmountToAdmin failed:%s", err)
		return false
	}
	balance, err := e.GetBalance()
	if err != nil {
		return false
	}

	return balance[0].Amount >= int64(value)
}

type ApplyRequest struct {
	Amount int64  `json:"amount"`
	Memo   string `json:"memo"`
}

type PayInvoiceRequest struct {
	Invoice  string `json:"invoice"`
	FeeLimit int64  `json:"feeLimit"`
}

type PaymentRequest struct {
	AssetId string `json:"asset_id"`
}
