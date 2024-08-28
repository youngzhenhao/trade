package btc_channel

import (
	"trade/models"
	rpc "trade/services/servicesrpc"
)

var (
	ChannelBtcServiceFee = uint64(100)
	AssetServiceFee      = 100
)

func PayServerFee(account *models.Account, fee uint64) error {
	acc, err := rpc.AccountInfo(account.UserAccountCode)
	if err != nil {
		return err
	}
	// Change the escrow account balance
	_, err = rpc.AccountUpdate(account.UserAccountCode, acc.CurrentBalance-int64(fee), -1)
	if err != nil {
		return err
	}
	return nil
}
