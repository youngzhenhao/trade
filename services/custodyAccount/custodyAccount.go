package custodyAccount

import (
	"trade/services/btldb"
	"trade/services/custodyAccount/btc_channel"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/lockPayment"
)

func GetAssetBalanceList(userName string) *[]cBase.Balance {
	list := make(map[string]*cBase.Balance)

	e, err := btc_channel.NewBtcChannelEvent(userName)
	if err != nil {
		return nil
	}

	err, unlockedBalance, lockedBalance := lockPayment.GetBalance(e.UserInfo.User.Username, "00")
	if err != nil {
		return nil
	}
	list["00"] = &cBase.Balance{
		AssetId: "00",
		Amount:  int64(unlockedBalance + lockedBalance),
	}
	temp, err := btldb.GetAccountBalanceByAccountId(e.UserInfo.Account.ID)
	if err != nil {
		return nil
	}
	for _, v := range *temp {
		_, exists := list[v.AssetId]
		if exists {
			list[v.AssetId].Amount += int64(v.Amount)
		} else {
			list[v.AssetId] = &cBase.Balance{
				AssetId: v.AssetId,
				Amount:  int64(v.Amount),
			}
		}
	}
	getBalances, err := lockPayment.GetBalances(e.UserInfo.User.Username)
	if err != nil {
		return nil
	}
	if getBalances != nil {
		for _, v := range *getBalances {
			_, exists := list[v.AssetId]
			if exists {
				list[v.AssetId].Amount += int64(v.Amount)
			} else {
				list[v.AssetId] = &cBase.Balance{
					AssetId: v.AssetId,
					Amount:  int64(v.Amount),
				}
			}
		}
	}
	var result []cBase.Balance
	for _, v := range list {
		result = append(result, *v)
	}
	return &result
}
