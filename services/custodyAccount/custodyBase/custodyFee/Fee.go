package custodyFee

import (
	"trade/btlLog"
	"trade/models"
	"trade/services/btldb"
	"trade/services/custodyAccount/account"
	"trade/services/servicesrpc"
)

var (
	ChannelBtcInsideServiceFee = uint64(10)
	ChannelBtcServiceFee       = uint64(100)
	AssetInsideFee             = uint64(10)
	AssetOutsideFee            = uint64(2500)
)

func PayServerFee(u *account.UserInfo, fee uint64) error {
	acc, err := servicesrpc.AccountInfo(u.Account.UserAccountCode)
	if err != nil {
		return err
	}
	// Change the escrow account balance
	_, err = servicesrpc.AccountUpdate(u.Account.UserAccountCode, acc.CurrentBalance-int64(fee), -1)
	if err != nil {
		return err
	}
	return nil
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
