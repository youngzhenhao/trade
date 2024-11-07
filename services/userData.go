package services

import (
	"errors"
	"gopkg.in/yaml.v3"
	"time"
	"trade/btlLog"
	"trade/middleware"
	"trade/models"
	"trade/models/custodyModels"
	"trade/services/btldb"
	"trade/services/custodyAccount/btc_channel"
	"trade/services/custodyAccount/custodyBase"
	"trade/utils"
)

func ReadLockBalanceByAccountId(accountId int) (*[]custodyModels.LockBalance, error) {
	var lockBalances []custodyModels.LockBalance
	err := middleware.DB.Where("account_id = ?", accountId).Find(&lockBalances).Error
	return &lockBalances, err
}

func ReadAccountBalanceByAccountId(accountId int) (*[]models.AccountBalance, error) {
	var accountBalances []models.AccountBalance
	err := middleware.DB.Where("account_id = ?", accountId).Find(&accountBalances).Error
	return &accountBalances, err
}

func GetUserInfoData(username string) (*models.UserInfoData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	user, err := btldb.ReadUserByUsername(username)
	if err != nil {
		return nil, utils.AppendErrorInfo(err, "ReadUserByUsername")
	}
	account, err := btldb.ReadAccountByName(username)
	if err != nil {
		btlLog.UserData.Error("ReadAccountByName err:%v", err)
		account = &models.Account{}
	}
	userInfoData := models.UserInfoData{
		Username:     username,
		UserId:       int(user.ID),
		CreatedAt:    user.CreatedAt,
		Account:      account.UserAccountCode,
		AccountId:    int(account.ID),
		UserRecentIp: user.RecentIpAddresses,
	}
	if err != nil {
		// @dev: return err here
		return &userInfoData, err
	}
	return &userInfoData, nil
}

func GetUserBtcBalanceData(username string) (*models.UserBtcBalanceData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	btcBalance, err := GetBtcBalanceByUsername(username)
	if err != nil {
		btlLog.UserData.Error("GetBtcBalanceByUsername err:%v", err)
		btcBalance = &models.BtcBalance{}
	}
	userBtcBalanceData := models.UserBtcBalanceData{
		CreatedAt:          btcBalance.CreatedAt,
		UpdatedAt:          btcBalance.UpdatedAt,
		TotalBalance:       btcBalance.TotalBalance,
		ConfirmedBalance:   btcBalance.ConfirmedBalance,
		UnconfirmedBalance: btcBalance.UnconfirmedBalance,
		LockedBalance:      btcBalance.LockedBalance,
	}
	return &userBtcBalanceData, nil
}

func GetUserAssetBalanceData(username string) (*[]models.UserAssetBalanceData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	assetBalances, err := btldb.ReadAssetBalanceByUsername(username)
	if err != nil {
		btlLog.UserData.Error("ReadAssetBalanceByUsername err:%v", err)
		assetBalances = &[]models.AssetBalance{}
	}
	var userAssetBalanceDatas []models.UserAssetBalanceData
	for _, assetBalance := range *assetBalances {
		userAssetBalanceDatas = append(userAssetBalanceDatas, models.UserAssetBalanceData{
			CreatedAt: assetBalance.CreatedAt,
			UpdatedAt: assetBalance.UpdatedAt,
			AssetId:   assetBalance.AssetID,
			Name:      assetBalance.Name,
			AssetType: assetBalance.AssetType,
			Balance:   assetBalance.Balance,
		})
	}
	return &userAssetBalanceDatas, nil
}

func GetUserAddrReceiveData(username string) (*[]models.UserAddrReceiveData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	addrReceives, err := btldb.ReadAddrReceiveEventsByUsername(username)
	if err != nil {
		btlLog.UserData.Error("ReadAddrReceiveEventsByUsername err:%v", err)
		addrReceives = &[]models.AddrReceiveEvent{}
	}
	var userAddrReceiveDatas []models.UserAddrReceiveData
	for _, addrReceive := range *addrReceives {
		userAddrReceiveDatas = append(userAddrReceiveDatas, models.UserAddrReceiveData{
			CreatedAt:    addrReceive.CreatedAt,
			UpdatedAt:    addrReceive.UpdatedAt,
			CreationTime: utils.TimestampToTime(addrReceive.CreationTimeUnixSeconds),
			AssetId:      addrReceive.AddrAssetID,
			AddrEncoded:  addrReceive.AddrEncoded,
			Amount:       addrReceive.AddrAmount,
			Outpoint:     addrReceive.Outpoint,
		})
	}
	return &userAddrReceiveDatas, nil
}

func GetUserAssetTransferData(username string) (*[]models.UserAssetTransferData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	assetTransfers, err := btldb.ReadAssetTransferProcessedSliceByUsername(username)
	if err != nil {
		btlLog.UserData.Error("ReadAssetTransferProcessedSliceByUsername err:%v", err)
		assetTransfers = &[]models.AssetTransferProcessedDb{}
	}
	var userAssetTransferDatas []models.UserAssetTransferData
	for _, assetTransfer := range *assetTransfers {
		userAssetTransferDatas = append(userAssetTransferDatas, models.UserAssetTransferData{
			CreatedAt:    assetTransfer.CreatedAt,
			UpdatedAt:    assetTransfer.UpdatedAt,
			Txid:         assetTransfer.Txid,
			AssetId:      assetTransfer.AssetID,
			TransferTime: utils.TimestampToTime(assetTransfer.TransferTimestamp),
		})
	}
	return &userAssetTransferDatas, nil
}

func GetUserAccountBtcBalanceData(username string) (*models.UserAccountBtcBalanceData, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	var balance int64
	btcChannelEvent, err := btc_channel.NewBtcChannelEvent(username)
	if err != nil {
		// TODO: return error info
		btlLog.UserData.Error("NewBtcChannelEvent err:%v", err)
		btcChannelEvent = &btc_channel.BtcChannelEvent{}
	} else {
		getBalance, err := btcChannelEvent.GetBalance()
		if err != nil {
			btlLog.UserData.Error("NewBtcChannelEvent err:%v", err)
			getBalance = []custodyBase.Balance{}
		}
		balance = getBalance[0].Amount
	}
	var userAccountBtcBalanceData models.UserAccountBtcBalanceData
	userAccountBtcBalanceData = models.UserAccountBtcBalanceData{
		Amount: int(balance),
	}
	return &userAccountBtcBalanceData, nil
}

func GetUserAccountAssetBalanceData(accountId int) (*[]models.UserAccountAssetBalanceData, error) {
	if accountId == 0 {
		return nil, errors.New("account id is zero")
	}
	accountBalances, err := ReadAccountBalanceByAccountId(accountId)
	var userAccountAssetBalanceDatas []models.UserAccountAssetBalanceData
	if err != nil {
		btlLog.UserData.Error("ReadLockBalanceByAccountId err:%v", err)
		accountBalances = &[]models.AccountBalance{}
	}
	for _, accountBalance := range *accountBalances {
		userAccountAssetBalanceDatas = append(userAccountAssetBalanceDatas, models.UserAccountAssetBalanceData{
			CreatedAt: accountBalance.CreatedAt,
			UpdatedAt: accountBalance.UpdatedAt,
			AssetId:   accountBalance.AssetId,
			Amount:    accountBalance.Amount,
		})
	}
	return &userAccountAssetBalanceDatas, nil
}

// GetUserData
// @Description: Get user data
func GetUserData(username string) (*models.UserData, error) {
	var userData models.UserData
	errorInfos := new([]string)
	userInfo, err := GetUserInfoData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserInfoData")
	}
	userBtcBalance, err := GetUserBtcBalanceData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserBtcBalanceData")
	}
	userAssetBalance, err := GetUserAssetBalanceData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserAssetBalanceData")
	}
	userAddrReceive, err := GetUserAddrReceiveData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserAddrReceiveData")
	}
	userAssetTransfer, err := GetUserAssetTransferData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserAssetTransferData")
	}
	userAccountBtcBalance, err := GetUserAccountBtcBalanceData(username)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserAccountBtcBalanceData")
	}
	var accountId int
	accountId = userInfo.AccountId
	userAccountAssetBalance, err := GetUserAccountAssetBalanceData(accountId)
	if err != nil {
		*errorInfos = append(*errorInfos, err.Error())
		//return nil, utils.AppendErrorInfo(err, "GetUserAccountAssetBalanceData")
	}
	userData = models.UserData{
		QueryTime:               time.Now(),
		UserInfo:                userInfo,
		UserBtcBalance:          userBtcBalance,
		UserAssetBalance:        userAssetBalance,
		UserAddrReceive:         userAddrReceive,
		UserAssetTransfer:       userAssetTransfer,
		UserAccountBtcBalance:   userAccountBtcBalance,
		UserAccountAssetBalance: userAccountAssetBalance,
		ErrorInfos:              errorInfos,
	}
	return &userData, nil
}

func GetUserDataYaml(username string) (string, error) {
	userData, err := GetUserData(username)
	if err != nil {
		return "", utils.AppendErrorInfo(err, "GetUserData")
	}
	userDataBytes, _ := yaml.Marshal(userData)
	return string(userDataBytes), nil
}
