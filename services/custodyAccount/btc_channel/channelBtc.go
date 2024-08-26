package btc_channel

import (
	"trade/btlLog"
	caccount "trade/services/custodyAccount/account"
	"trade/services/custodyAccount/custodyBase"
	rpc "trade/services/servicesrpc"
)

type BtcChannelEvent struct {
	custodyBase.CustodyEvent
	UserInfo *caccount.UserInfo
}

func NewBtcChannelEvent(UserName string) (*BtcChannelEvent, error) {
	var (
		e   BtcChannelEvent
		err error
	)
	e.UserInfo, err = caccount.GetUserInfo(UserName)
	if err != nil {
		btlLog.CUST.Error(err.Error())
		return nil, err
	}
	return &e, nil
}

func (e *BtcChannelEvent) GetBalance() (*[]custodyBase.Balance, error) {
	acc, err := rpc.AccountInfo(e.UserInfo.Account.UserAccountCode)
	if err != nil {
		return nil, err
	}
	balances := []custodyBase.Balance{
		{
			AssetId: "00",
			Amount:  acc.CurrentBalance,
		},
	}
	return &balances, nil
}
