package custodyAccount

import (
	"sync"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/btc_channel"
)

// CreateBackFeeMission 创建退费任务
func CreateBackFeeMission(payInsideId uint) (uint, error) {
	//判断是否存在支付记录
	_, err := btldb.ReadPayInside(payInsideId)
	if err != nil {
		btlLog.CUST.Error("CreateBackFeeMission find payInside error:%v", err.Error())
		return 0, err
	}
	//创建退费任务
	var backMission models.BackFee
	backMission.PayInsideId = payInsideId
	backMission.Status = models.BackFeeStatePending
	err = middleware.DB.Create(&backMission).Error
	if err != nil {
		btlLog.CUST.Error("CreateBackFeeMission error:%v", err.Error())
		return 0, err
	}
	return backMission.ID, nil
}

var BackFeeMutex sync.Mutex

// PollBackFeeMission 处理退费任务
func PollBackFeeMission() {
	BackFeeMutex.Lock()
	defer BackFeeMutex.Unlock()
	var results []BackFeeSqlResult
	middleware.DB.Raw(getBackFeeSql, 0).Scan(&results)
	for _, r := range results {
		account, err := btldb.ReadAccountByUserId(r.PayUserId)
		if err != nil {
			btlLog.CUST.Error("PollBackFeeMission find pay account error:%v", err.Error())
			continue
		}
		balanceId, err := btc_channel.UpdateCustodyAccount(account, models.AWAY_IN, r.GasFee+r.ServeFee, "backFee", 0)
		if err != nil {
			btlLog.CUST.Error("PollBackFeeMission update custody account error:%v", err.Error())
			continue
		}
		err = middleware.DB.Table("user_back_fees").Where("user_back_fees.id = ?", r.BackFeeId).Update("back_balance_id", balanceId).Update("status", models.BackFeeStatePaid).Error
		if err != nil {
			btlLog.CUST.Error("PollBackFeeMission update balance_id error:%v", err.Error())
			continue
		}
	}
}

type BackFeeSqlResult struct {
	BackFeeId   uint
	PayInsideId uint
	PayUserId   uint
	GasFee      uint64
	ServeFee    uint64
	PayType     string
	AssetType   string
}

const getBackFeeSql = `
	SELECT
		user_back_fees.id AS BackFeeId,
		user_back_fees.pay_inside_id AS PayInsideId,
		user_pay_inside.pay_user_id AS  PayUserId,
		user_pay_inside.gas_fee AS GasFee,
		user_pay_inside.serve_fee AS ServeFee,
		user_pay_inside.pay_type AS PayType,
		user_pay_inside.asset_type AS AssetType
	FROM 
	    user_back_fees
	LEFT JOIN 
	    user_pay_inside ON user_back_fees.pay_inside_id = user_pay_inside.id
	WHERE 
	    user_back_fees.status = ?
`

// checkBackFeeMissionById 检查退费任务状态是否成功
func checkBackFeeMissionById(BackFeeId uint) bool {
	var BackFee models.BackFee
	err := middleware.DB.Where("id = ?", BackFeeId).First(&BackFee).Error
	if err != nil {
		btlLog.CUST.Error("CheckBackFeeMissionById find backFee error:%v", err.Error())
		return false
	}
	if BackFee.Status == models.BackFeeStatePaid {
		return true
	}
	return false
}
