package custodyFee

import (
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase/custodyRpc"
)

var (
	ChannelBtcInsideServiceFee = uint64(10)
	ChannelBtcServiceFee       = uint64(100)
	AssetInsideFee             = uint64(10)
	AssetOutsideFee            = uint64(2500)
)

func PayServerFee(u *account.UserInfo, fee uint64) error {
	_, err := custodyRpc.UpdateBalance(u, custodyRpc.UpdateBalanceMinus, int64(fee))
	return err
}

func PayServiceFeeSync(u *account.UserInfo, fee uint64, balanceId uint, PayType models.PayInsideType, memo string) error {
	switch PayType {
	case models.ChannelBTCFee:
	default:
	}
	//划扣手续费
	err := PayServerFee(u, fee)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return err
	}
	//创建内部转账任务
	payInside := models.PayInside{
		PayUserId:     u.User.ID,
		GasFee:        fee,
		ServeFee:      0,
		ReceiveUserId: 1,
		PayType:       PayType,
		AssetType:     "00",
		PayReq:        &memo,
		BalanceId:     balanceId,
		Status:        models.PayInsideStatusSuccess,
	}
	//写入数据库
	err = btldb.CreatePayInside(&payInside)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return err
	}
	return err
}

func SetMemoSign() string {
	return "internal transfer"
}
