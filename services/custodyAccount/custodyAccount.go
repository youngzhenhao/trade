package custodyAccount

import (
	"trade/middleware"
	cBase "trade/services/custodyAccount/custodyBase"
	"trade/services/custodyAccount/defaultAccount/custodyAssets"
	"trade/services/custodyAccount/defaultAccount/custodyBtc"
	"trade/services/custodyAccount/lockPayment"
)

func GetAssetBalanceList(userName string) (*[]cBase.Balance, error) {
	list := make(map[string]*cBase.Balance)
	e, err := custodyBtc.NewBtcChannelEvent(userName)
	if err != nil {
		return nil, err
	}
	err, unlockedBalance, lockedBalance, _ := lockPayment.GetBalance(e.UserInfo.User.Username, "00")
	if err != nil {
		return nil, err
	}
	list["00"] = &cBase.Balance{
		AssetId: "00",
		Amount:  int64(unlockedBalance + lockedBalance),
	}
	temp := custodyAssets.GetAssetsBalances(middleware.DB, e.UserInfo.Account.ID)
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
		return nil, err
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
	return &result, nil
}
